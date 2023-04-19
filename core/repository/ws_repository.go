package repository

import (
	"context"
	"database/sql"
	"log"
	"soporte-go/core/model"
	"soporte-go/core/model/ws"
	"time"

	"soporte-go/core/model/caso"

	// "github.com/lib/pq"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type wsRepository struct {
	Conn    *sql.DB
	Context context.Context
}

func NewWsRepository(conn *sql.DB, ctx context.Context) ws.WsRepository {
	return &wsRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func (p *wsRepository) GetMessages(ctx context.Context, casoId string) (res []ws.Message, err error) {
	query := `select * from messages where caso_id = $1;`
	res, err = p.fetchMessages(ctx, query, casoId)

	return
}

func (p *wsRepository) SaveMessage(ctx context.Context, m *ws.MessageData) (res ws.Message, err error) {
	var message ws.Message
	var query string
	var caso caso.Caso
	query = `select estado,client_id,funcionario_id from casos where id = $1;`
	err = p.Conn.QueryRowContext(ctx, query, m.CasoId).Scan(&caso.Estado, &caso.ClienteId, &caso.FuncionarioId)
	if err != nil {
		log.Println(err)
	}
	if *caso.ClienteId == m.FromUser {
		query = `update casos set estado = $1,updated_on = $2 where id = $3;`
		_, err := p.Conn.ExecContext(ctx, query, model.EnEsperaDelFuncionario, time.Now(), m.CasoId)
		if err != nil {
			log.Println(err)
		}
	}
	if *caso.FuncionarioId == m.FromUser {
		query = `update casos set estado = $1,updated_on = $2 where id = $3;`
		_, err := p.Conn.ExecContext(ctx, query, model.EnEsperaDelCliente, time.Now(), m.CasoId)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println(m.MediaUrl)
	query = `insert into messages (from_user,to_user,caso_id,media_url,content,is_read,created_on) 
	values ($1,$2,$3,$4,$5,$6,$7) returning (id,from_user,to_user,caso_id,media_url,content,is_read,created_on,is_deleted);`
	err = p.Conn.QueryRowContext(ctx, query, m.FromUser, m.ToUser, m.CasoId,pq.Array(m.MediaUrl), m.Content, m.IsRead, time.Now()).Scan(&message)
	return
}

func (p *wsRepository) fetchMessages(ctx context.Context, query string, args ...interface{}) (result []ws.Message, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]ws.Message, 0)
	for rows.Next() {
		t := ws.Message{}
		err = rows.Scan(
			&t.Id,
			&t.CasoId,
			&t.FromUser,
			&t.ToUser,
			pq.Array(&t.MediaUrl),
			&t.Content,
			&t.IsRead,
			&t.CreatedOn,
			&t.IsDeleted,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}
