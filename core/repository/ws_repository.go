package repository

import (
	"context"
	"soporte-go/core/model/ws"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type wsRepository struct {
	Conn    *pgxpool.Pool
	Context context.Context
}

func NewWsRepository(conn *pgxpool.Pool, ctx context.Context) ws.WsRepository {
	return &wsRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func (p *wsRepository) GetMessages(ctx context.Context, casoId string)(res []ws.Message,err error){
	query := `select * from messages where caso_id = $1;`
	res,err = p.fetchMessages(ctx,query,casoId)
	return
}

func (p *wsRepository) SaveMessage(ctx context.Context,m *ws.MessageData)(res ws.Message,err error){
	query := `insert into messages (client_id,funcionario_id,caso_id,client_name,funcionario_name,media_url,content,is_read,created_on) 
	values ($1,$2,$3,$4,$5,$6,$7,$8,$9);`
	_,err = p.Conn.Exec(ctx,query,m.ClienteId,m.FuncionarioId,m.CasoId,m.ClienteName,m.FuncionarioName,m.MediaUrl,m.Content,m.IsRead,time.Now())
	return
}

func (p *wsRepository) fetchMessages(ctx context.Context, query string, args ...interface{}) (result []ws.Message, err error) {
	rows, err := p.Conn.Query(p.Context, query, args...)
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
			&t.ClienteId,
			&t.FuncionarioId,
			&t.ClienteName,
			&t.FuncionarioName,
			&t.MediaUrl,
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
