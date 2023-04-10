package caso

import (
	"context"
	"database/sql"
	"time"

	// "database/sql"
	"log"
	// "soporte-go/core/model"
	"soporte-go/core/model/caso"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type pgCasoRepository struct {
	Conn    *sql.DB
	Context context.Context
}

type ProfileInfo struct {
	EmpresaId int   `json:"empresa"`
	Areas     []int `json:"areas"`
}

func NewPgCasoRepository(conn *sql.DB, ctx context.Context) caso.CasoRepository {
	return &pgCasoRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func (p *pgCasoRepository) AsignarFuncionario(ctx context.Context, id string, idF string) (err error) {
	var query string
	var fechaInicio *string
	query = `update casos set funcionario_id = $1 where id = $2 returning (fecha_inicio)`
	err = p.Conn.QueryRowContext(ctx, query, idF, id).Scan(&fechaInicio)
	log.Println(fechaInicio)
	if err != nil {
		return
	}
	if fechaInicio == nil {
		query = `update casos set fecha_inicio = $1 where id = $2`
		_,err = p.Conn.ExecContext(ctx,query,time.Now(),id)

	}
	// query = `update casos set fecha_iniciap`
	return err
}

func (p *pgCasoRepository) GetCaso(ctx context.Context, id string) (res caso.Caso, err error) {
	log.Println(id)
	query := `select clientes.nombre,funcionarios.nombre,titulo,id,descripcion,detalles_de_finalizacion,empresa,area,casos.created_on,
	casos.updated_on,fecha_inicio,fecha_fin,prioridad,casos.estado,casos.client_id,casos.funcionario_id,casos.superior_id
	from casos inner join clientes on clientes.client_id = casos.client_id left join funcionarios on funcionarios.funcionario_id = casos.funcionario_id
	where id = $1`
	var casoDetail caso.Caso
	err = p.Conn.QueryRowContext(ctx, query, id).Scan(
		&casoDetail.ClienteName,
		&casoDetail.FuncionarioName,
		&casoDetail.Titulo,
		&casoDetail.Id,
		&casoDetail.Descripcion,
		&casoDetail.DetallesDeFinalizacion,
		&casoDetail.Empresa,
		&casoDetail.Area,
		&casoDetail.CreatedOn,
		&casoDetail.UpdatedOn,
		&casoDetail.FechaInicio,
		&casoDetail.FechaFin,
		&casoDetail.Prioridad,
		&casoDetail.Estado,
		&casoDetail.ClienteId,
		&casoDetail.FuncionarioId,
		&casoDetail.SuperiorId,
	)
	if err != nil {
		return caso.Caso{}, err
	}
	return casoDetail, err
}

func (p *pgCasoRepository) GetCasosCountFuncionario(ctx context.Context, id string) (res int, err error) {
	query3 := `select count(*) from casos where funcionario_id = $1`
	err = p.Conn.QueryRowContext(ctx, query3, id).Scan(&res)
	return
}

func (p *pgCasoRepository) GetCasosCountCliente(ctx context.Context, id string) (res int, err error) {
	query3 := `select count(*) from casos where client_id = $1`
	err = p.Conn.QueryRowContext(ctx, query3, id).Scan(&res)
	return
}

func (p *pgCasoRepository) GetCasosCountbySuperiorId(ctx context.Context, id string) (res int, err error) {
	query3 := `select count(*) from casos where superior_id = $1`
	err = p.Conn.QueryRowContext(ctx, query3, id).Scan(&res)
	return
}

func (p *pgCasoRepository) GetCasosCount(ctx context.Context) (res int, err error) {
	query3 := `select count(*) from casos`
	err = p.Conn.QueryRowContext(ctx, query3).Scan(&res)
	return
}

func (p *pgCasoRepository) GetCasosCliente(ctx context.Context, id string, query *caso.CasoQuery) (list []caso.Caso, err error) {
	// var count int
	query2 := `select id,titulo,created_on,updated_on,prioridad,estado
		from casos where client_id = $1 limit $2 offset $3`
	list, _ = p.fetchCasos(ctx, query2, id, query.PageSize, query.Page*10)
	// query3 := `select count(*) from casos where client_id = $1;`
	// err = p.Conn.QueryRowContext(ctx,query3,id).Scan(&count)
	// size = count
	return
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetCasosFuncionario(ctx context.Context, id string, query *caso.CasoQuery) (list []caso.Caso, err error) {
	query2 := `select id,titulo,created_on,updated_on,prioridad,estado
		from casos where funcionario_id = $1 limit $2 offset $3`
	list, _ = p.fetchCasos(ctx, query2, id, query.PageSize, query.Page*10)

	return
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetAllCasosUserCliente(ctx context.Context, id string, query *caso.CasoQuery) (list []caso.Caso, err error) {
	// var superiorId string
	// query1 := `select superior_id from clientes where client_id = $1;`
	// if err = p.Conn.QueryRowContext(ctx, query1, id).Scan(&superiorId); err != nil {
	// 	return
	// }
	query2 := `select id,titulo,created_on,updated_on,prioridad,estado
	 from casos where superior_id = $1 limit 10 offset $2;`
	if query.Page == 1 || query.Page == 0 {
		list, err = p.fetchCasos(ctx, query2, id, 0)
	} else {
		page := query.Page - 1
		list, err = p.fetchCasos(ctx, query2, id, page*10)
	}
	return
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetAllCasosUserFuncionario(ctx context.Context,id int,query *caso.CasoQuery)(res []caso.Caso,err error){
	// log.Println("Get all casos funcionario")
	query2 := `select id,titulo,created_on,updated_on,prioridad,estado
		from casos limit $1 offset $2`
	res, err = p.fetchCasos(ctx, query2,query.PageSize, query.Page*10)
	return
}

func (p *pgCasoRepository) UpdateCaso(ctx context.Context, columns []string, values ...interface{}) error {
	return nil
}

func (p *pgCasoRepository) StoreCaso(ctx context.Context, cas *caso.Caso, id string, emI int) (idCaso string, err error) {
	var superiorId string
	query := `select superior_id from clientes where client_id = $1;`
	err = p.Conn.QueryRowContext(ctx, query, id).Scan(&superiorId)
	if err != nil {
		return
	}
	var casoId string
	// query := `select (empresa_id) from clientes where user_id = $1;`
	// err = p.Conn.QueryRowContext(ctx, query, id).Scan(&empresaId)
	// if err != nil {
	// log.Println(err)
	// }
	// log.Println(empresaId)
	query1 := `insert into casos (titulo,client_id,descripcion,empresa,area,prioridad,created_on,superior_id) values(
		$1,$2,$3,$4,$5,$6,$7,$8) returning (id);`
	err = p.Conn.QueryRowContext(ctx, query1, cas.Titulo, id, cas.Descripcion, emI, cas.Area,
		cas.Prioridad, time.Now(), superiorId).Scan(&casoId)
	if err != nil {
		log.Println(err)
	}
	return casoId, err
}

func (p *pgCasoRepository) fetchCasos(ctx context.Context, query string, args ...interface{}) (result []caso.Caso, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
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
			&t.CreatedOn,
			&t.UpdatedOn,
			&t.Prioridad,
			&t.Estado,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}
