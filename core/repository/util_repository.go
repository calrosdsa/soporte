package repository

import (
	"context"
	"database/sql"
	"log"
	"soporte-go/core/model/util"
)



type pgUtilRepository struct {
	Conn *sql.DB
}

func NewPgUtilRepository(c *sql.DB) util.UtilRepository {
	return &pgUtilRepository{
		Conn: c,
	}
}

func (p *pgUtilRepository) GetEmailFromUserCasos(ctx context.Context,id int,userId string)(res util.NotificationCasoEmail,err error){
	// var mail string
	var query string
	// var d util.NotificationCasoEmail
	query = `select f.email from funcionarios as f 
	inner join user_area as ua on  f.funcionario_id = ua.user_id where ua.area_id = $1`
	mails,err := p.fetchEmails(ctx,query,id)
	if err != nil {
		log.Println(err)
	}
	query = `select c.nombre,c.apellido,p.nombre,e.nombre from clientes as c inner join proyectos as p on p.id = $2
	inner join empresas as e on e.id = c.empresa_id where client_id = $1`
	err = p.Conn.QueryRowContext(ctx,query,userId,id).Scan(
		&res.UsuarioName,
		&res.UsuarioApellido,
		&res.Area,
		&res.Entidad,
	)
	if err != nil {
		log.Println(err)
	}
	res.Mails = mails
	log.Println(id)
	log.Println("mails===--",mails)
	return
}

func (p *pgUtilRepository) fetchEmails(ctx context.Context, query string, args ...interface{}) (result []util.Mails, err error) {
	rows, err := p.Conn.QueryContext(ctx , query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]util.Mails, 0)
	for rows.Next() {
		t := util.Mails{}
		err = rows.Scan(
			&t.Mail,
		)
		result = append(result, t)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return result, nil
}