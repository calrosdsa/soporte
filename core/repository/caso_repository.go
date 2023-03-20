package repository

import (
	"context"
	"time"

	// "database/sql"
	"log"
	"soporte-go/core/model"
	"soporte-go/core/model/caso"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type pgCasoRepository struct {
	Conn    *pgxpool.Pool
	Context context.Context
}

type ProfileInfo struct {
	EmpresaId int   `json:"empresa"`
	Areas     []int `json:"areas"`
}

func NewPgCasoRepository(conn *pgxpool.Pool, ctx context.Context) caso.CasoRepository {
	return &pgCasoRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func (p *pgCasoRepository) AsignarFuncionario(ctx context.Context,id string,idF string)(err error){
	query := `update casos set funcionario_id = $1,funcionario_name = 'Jorge M' where id = $2;`
	_,err = p.Conn.Exec(ctx,query,idF,id)
	return err
}	

func (p *pgCasoRepository) GetCaso(ctx context.Context, id string) (res caso.Caso, err error) {
	log.Println(id)
	query := `select * from casos where id = $1;`
	list, err := p.fetchCasos(ctx, query, id)
	if err != nil {
		return caso.Caso{}, err
	}
	if len(list) > 0 {
		res = list[0]
	} else {
		return res, model.ErrNotFound
	}
	return
}

func (p *pgCasoRepository) GetCasosCliente(ctx context.Context, id *string,query *caso.CasoQuery) (list []caso.Caso,size int,err error) {
	var count int	
		query2 := `select * from casos where client_id = $1 limit 10 offset $2;`
		if query.Page == 1 || query.Page == 0 {
		list, _ = p.fetchCasos(ctx, query2, id,0)
		}else{
			page := query.Page - 1
			list, _ = p.fetchCasos(ctx, query2, id,page * 10)
		}
		query3 := `select count(*) from casos where client_id = $1;`
		err = p.Conn.QueryRow(ctx,query3,id).Scan(&count)
		size = count
	return 
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetCasosFuncionario(ctx context.Context, id *string,query *caso.CasoQuery) (list []caso.Caso,size int,err error) {
	var count int	
		query2 := `select * from casos where funcionario_id = $1 limit 10 offset $2;`
		if query.Page == 1 || query.Page == 0 {
		list, _ = p.fetchCasos(ctx, query2, id,0)
		}else{
			page := query.Page - 1
			list, _ = p.fetchCasos(ctx, query2, id,page * 10)
		}
		query3 := `select count(*) from casos where client_id = $1;`
		err = p.Conn.QueryRow(ctx,query3,id).Scan(&count)
		size = count
	return 
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetAllCasosUser(ctx context.Context, id string,query *caso.CasoQuery) (list []caso.Caso,size int,err error) {
	var superiorId string;
	query1 := `select superior_id from clientes where client_id = $1;`
	err = p.Conn.QueryRow(ctx,query1,id).Scan(&superiorId)

	var count int
	query2 := `select * from casos where superior_id = $1 limit 10 offset $2;`
	if query.Page == 1 || query.Page == 0 {
		list, err = p.fetchCasos(ctx, query2, superiorId,0)
	}else{
		page := query.Page - 1
		list, err = p.fetchCasos(ctx, query2, id,page * 10)
	}
	query3 := `select count(*) from casos where superior_id = $1;`
	err = p.Conn.QueryRow(ctx,query3,superiorId).Scan(&count)
	size = count
	return 
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) UpdateCaso(ctx context.Context, columns []string, values ...interface{}) error {
	return nil
}

func (p *pgCasoRepository) StoreCaso(ctx context.Context, cas *caso.Caso, id string,emI int) (idCaso string, err error) {
	var superiorId string;
	query := `select superior_id from clientes where client_id = $1;`
	err = p.Conn.QueryRow(ctx,query,id).Scan(&superiorId)
	var casoId string
	// query := `select (empresa_id) from clientes where user_id = $1;`
	// err = p.Conn.QueryRow(ctx, query, id).Scan(&empresaId)
	// if err != nil {
		// log.Println(err)
	// }
	// log.Println(empresaId)
	query1 := `insert into casos (titulo,client_id,descripcion,empresa,area,cliente_name,prioridad,created_on,superior_id) values(
		$1,$2,$3,$4,$5,$6,$7,$8,$9) returning (id);`
	err = p.Conn.QueryRow(ctx, query1, cas.Titulo, id, cas.Descripcion, emI, cas.Area, cas.ClienteName,
		cas.Prioridad,time.Now(),superiorId).Scan(&casoId)
	if err != nil {
		log.Println(err)
	}
	return casoId, err
}

func (p *pgCasoRepository) fetchCasos(ctx context.Context, query string, args ...interface{}) (result []caso.Caso, err error) {
	rows, err := p.Conn.Query(p.Context, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]caso.Caso, 0)
	for rows.Next() {
		t := caso.Caso{}
		err = rows.Scan(
			&t.Id,
			&t.Titulo,
			&t.Descripcion,
			&t.DetallesDeFinalizacion,
			&t.Empresa,
			&t.Area,
			&t.ClienteName,
			&t.FuncionarioName,
			&t.CreatedOn,
			&t.UpdatedOn,
			&t.FechaInicio,
			&t.FechaFin,
			&t.Prioridad,
			&t.Estado,
			&t.ClienteId,
			&t.FuncionarioId,
			&t.SuperiorId,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}
