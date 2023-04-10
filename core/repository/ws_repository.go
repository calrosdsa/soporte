package repository

import (
	"context"
	"log"
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
	var message ws.Message;
	query := `insert into messages (from_user,to_user,caso_id,media_url,content,is_read,created_on) 
	values ($1,$2,$3,$4,$5,$6,$7) returning (id,from_user,to_user,caso_id,media_url,content,is_read,created_on,is_deleted);`
	err = p.Conn.QueryRow(ctx,query,m.FromUser,m.ToUser,m.CasoId,m.MediaUrl,m.Content,m.IsRead,time.Now()).Scan(&message)
	log.Println(err)
	log.Println(message.Id)
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
			&t.FromUser,
			&t.ToUser,
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
