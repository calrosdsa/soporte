package caso

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// "database/sql"
	"log"
	// "soporte-go/core/model"
	"soporte-go/core/model"
	"soporte-go/core/model/caso"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
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

func (p *pgCasoRepository) GetCasosCliForReporte(ctx context.Context,options *caso.CasoReporteOptions)(res []caso.Caso,err error) {
	var query string
	log.Println(options.Estados)
	if len(options.Estados) == 3{
		// log.Println("ALL estaods")
		query = `select clientes.nombre,clientes.apellido,funcionarios.nombre,funcionarios.apellido,titulo,C.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
		c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.client_id,c.funcionario_id,c.superior_id,c.rol,p.nombre
		from casos as c inner join clientes on clientes.client_id = c.client_id left join funcionarios
		 on funcionarios.funcionario_id = c.funcionario_id
		left join proyectos as p on c.area = p.id
		 where c.created_on between $1 and $2 and c.area = any($3);`
		res,err = p.fetchCasosDetail(ctx,query,options.StartDate,options.EndDate,pq.Array(options.Areas))
	}else{
		query = `select clientes.nombre,clientes.apellido,funcionarios.nombre,funcionarios.apellido,titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
		c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.client_id,c.funcionario_id,c.superior_id,c.rol,p.nombre
		from casos as c inner join clientes on clientes.client_id = c.client_id 
		left join funcionarios on funcionarios.funcionario_id = c.funcionario_id
		left join proyectos as p on c.area = p.id
		 where c.created_on between $1 and $2 and c.estado = any($3) and c.area = any($4);`
		res,err = p.fetchCasosDetail(ctx,query,options.StartDate,options.EndDate,pq.Array(options.Estados),pq.Array(options.Areas))
	}
	return
}

func (p *pgCasoRepository) GetCasosFunForReporte(ctx context.Context,options *caso.CasoReporteOptions)(res []caso.Caso,err error){
	var query string
	log.Println(options.Estados)
	if len(options.Estados) == 3{
		// log.Println("ALL estaods")
		query = `select uc.nombre,uc.apellido,uf.nombre,uf.apellido,titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
		c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.funcionario_id,c.funcionario_id,c.superior_id,c.rol,p.nombre
		from casos as c inner join funcionarios as uc on uc.funcionario_id = c.client_id 
		left join funcionarios as uf on uf.funcionario_id = c.funcionario_id
		left join proyectos as p on c.area = p.id
		 where c.created_on between $1 and $2 and c.area = any($3);`
		res,err = p.fetchCasosDetail(ctx,query,options.StartDate,options.EndDate,pq.Array(options.Areas))
	}else{
		query = `select uc.nombre,uc.apellido,uf.nombre,uf.apellido,titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
		c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.funcionario_id,c.funcionario_id,c.superior_id,c.rol,p.nombre
		from casos as c inner join funcionarios as uc on uc.funcionario_id = c.client_id 
		left join funcionarios as uf on uf.funcionario_id = c.funcionario_id
		left join proyectos as p on c.area = p.id
		 where c.created_on between $1 and $2 and c.estado = any($3) and c.area = any($4);`
		res,err = p.fetchCasosDetail(ctx,query,options.StartDate,options.EndDate,pq.Array(options.Estados),pq.Array(options.Areas))
	}	
	return
}

func (p *pgCasoRepository) FinalizarCaso(ctx context.Context,fD *caso.FinalizacionDetail) (err error) {
	log.Println(fD.Detail)
	query := `update casos set detalles_de_finalizacion = $1, estado = $2,fecha_fin = $3,updated_on = $4 where id = $5`
	_,err = p.Conn.ExecContext(ctx,query,fD.Detail,fD.Estado,time.Now(),time.Now(),fD.Id)
	return
}

func (p *pgCasoRepository) AsignarFuncionario(ctx context.Context, id string, idF string) (err error) {
	var query string
	var fechaInicio *string
	query = `update casos set funcionario_id = $1,fecha_inicio = $2,updated_on = $3 where id = $4 returning (fecha_inicio)`
	err = p.Conn.QueryRowContext(ctx, query, idF, time.Now(),time.Now(),id).Scan(&fechaInicio)
	log.Println(fechaInicio)
	if err != nil {
		return
	}
	// if fechaInicio == nil {
	// 	query = `update casos set fecha_inicio = $1 where id = $2`
	// 	_, err = p.Conn.ExecContext(ctx, query, time.Now(), id)

	// }
	// query = `update casos set fecha_iniciap`
	return err
}



func (p *pgCasoRepository) GetCasoCliente(ctx context.Context, id string) (res caso.Caso, err error) {
	query := `select clientes.nombre,clientes.apellido,funcionarios.nombre,funcionarios.apellido,titulo,casos.id,descripcion,detalles_de_finalizacion,empresa,area,casos.created_on,
	casos.updated_on,fecha_inicio,fecha_fin,prioridad,casos.estado,casos.client_id,casos.funcionario_id,casos.superior_id,
	casos.rol,p.nombre
	from casos inner join clientes on clientes.client_id = casos.client_id 
	left join funcionarios on funcionarios.funcionario_id = casos.funcionario_id
	left join proyectos as p on casos.area = p.id
	where casos.id = $1 limit 1;`
	// var casoDetail caso.Caso
	list , err := p.fetchCasosDetail(ctx, query, id)
	if err != nil {
		return caso.Caso{}, err
	}
	res = list[0]
	return
}

func (p *pgCasoRepository) GetCasoFuncionario(ctx context.Context, id string) (res caso.Caso, err error) {
	query := `select 
	uc.nombre,uc.apellido,
	uf.nombre,uf.apellido,
	titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
	c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.client_id,c.funcionario_id,c.superior_id,c.rol,p.nombre
	from casos as c left join funcionarios as uc on uc.funcionario_id = c.client_id 
	left join proyectos as p on c.area = p.id
	left join funcionarios as uf on uf.funcionario_id = c.funcionario_id
	where c.id = $1 limit 1`
	// var casoDetail caso.Caso
	list , err := p.fetchCasosDetail(ctx, query, id)
	if err != nil {
		log.Println(err)
		return caso.Caso{}, err
	}
	res = list[0]
	return
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
	// query2 := `select id,titulo,created_on,updated_on,prioridad,estado,rol
		// from casos where client_id = $1 limit $2 offset $3`
	query2 := fmt.Sprintf(`select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	cl.nombre,cl.apellido,
	f2.nombre,f2.apellido
	from casos as c 
	left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	where c.client_id = $1 %s limit $2 offset $3`,query.Order)
	list, _ = p.fetchCasosWithClient(ctx, query2, id, query.PageSize, query.Page*10)
	// query3 := `select count(*) from casos where client_id = $1;`
	// err = p.Conn.QueryRowContext(ctx,query3,id).Scan(&count)
	// size = count
	return
	// return []caso.Caso{},nil
}

 


func (p *pgCasoRepository) GetCasosFuncionario(ctx context.Context, id string, query *caso.CasoQuery) (list []caso.Caso, err error) {
	query1 := fmt.Sprintf(`select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido,
	f2.nombre,f2.apellido
	from casos as c 
	left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f on f.funcionario_id = c.client_id
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	where c.funcionario_id = $1 or c.client_id = $2 %s limit $3 offset $4`,query.Order)
	list, err= p.fetchCasosWithClient(ctx, query1, id, id,query.PageSize, query.Page*10)
	return
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetAllCasosUserCliente(ctx context.Context, id string, query *caso.CasoQuery) (list []caso.Caso, err error) {
	// var superiorId string
	// query1 := `select superior_id from clientes where client_id = $1;`
	// if err = p.Conn.QueryRowContext(ctx, query1, id).Scan(&superiorId); err != nil {
	// 	return
	// }
	query2 := fmt.Sprintf(`select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	cl.nombre,cl.apellido,
	f2.nombre,f2.apellido
	from casos as c left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	where c.superior_id = $1 %s limit 10 offset $2`,query.Order)
	if query.Page == 1 || query.Page == 0 {
		list, err = p.fetchCasosWithClient(ctx, query2, id, 0)
	} else {
		page := query.Page - 1
		list, err = p.fetchCasosWithClient(ctx, query2, id, page*10)
	}
	return
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetAllCasosUserFuncionario(ctx context.Context, id int, query *caso.CasoQuery) (res []caso.Caso, err error) {
	// log.Println("Get all casos funcionario")
	log.Println(model.ASC)
	query2 := fmt.Sprintf(`select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido,
	f2.nombre,f2.apellido
	from casos as c left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	left join funcionarios as f on f.funcionario_id = c.client_id
	%s limit $1 offset $2`,query.Order)
	res, err = p.fetchCasosWithClient(ctx, query2, query.PageSize, query.Page*10)
	return
}

func (p *pgCasoRepository) UpdateCaso(ctx context.Context,c *caso.Caso) (err error) {
	query := `update casos set titulo = $1,descripcion = $2,updated_on = $3 where id = $4 
	returning titulo,descripcion,updated_on`
	err = p.Conn.QueryRowContext(ctx,query,c.Titulo,c.Descripcion,time.Now(),c.Id).Scan(
		&c.Titulo,&c.Descripcion,&c.UpdatedOn)
	if err != nil {
		log.Println(err)
	}
	return 
}

func (p *pgCasoRepository) CreateCasoCliente(ctx context.Context, cas *caso.Caso, id string, emI int,rol int) (err error) {
	query := `select superior_id from clientes where client_id = $1;`
	err = p.Conn.QueryRowContext(ctx, query, id).Scan(&cas.SuperiorId)
	if err != nil {
		return
	}
	query1 := `insert into casos (titulo,client_id,descripcion,empresa,area,prioridad,created_on,superior_id,rol) values(
		$1,$2,$3,$4,$5,$6,$7,$8,$9) returning id,created_on,rol;`
	err = p.Conn.QueryRowContext(ctx, query1, cas.Titulo, id, cas.Descripcion, emI, cas.Area,
		cas.Prioridad, time.Now(),cas.SuperiorId,rol).Scan(&cas.Id,&cas.CreatedOn,&cas.Rol)
	if err != nil {
		log.Println(err)
	}
	return
}

func (p *pgCasoRepository) CreateCasoFuncionario(ctx context.Context, cas *caso.Caso, id string, emI int,rol int) (err error) {
	log.Println("Create casp for funcionario")
	query := `select superior_id from funcionarios where funcionario_id = $1;`
	err = p.Conn.QueryRowContext(ctx, query, id).Scan(&cas.SuperiorId)
	if err != nil {
		return
	}
	if rol == int(model.RoleFuncionario) {
		query1 := `insert into casos (titulo,client_id,descripcion,empresa,area,prioridad,created_on,superior_id,funcionario_id,rol) values(
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id,created_on,rol;`
			err = p.Conn.QueryRowContext(ctx, query1, cas.Titulo, id, cas.Descripcion, emI, cas.Area,
				cas.Prioridad, time.Now(), cas.SuperiorId,cas.FuncionarioId,rol).Scan(&cas.Id,&cas.CreatedOn,&cas.Rol)
				if err != nil {
					log.Println(err)
				}
	} else {
		query1 := `insert into casos (titulo,client_id,descripcion,empresa,area,prioridad,created_on,superior_id,rol) values(
			$1,$2,$3,$4,$5,$6,$7,$8,$9) returning id,created_on,rol;`
			err = p.Conn.QueryRowContext(ctx, query1, cas.Titulo, id, cas.Descripcion, emI, cas.Area,
				cas.Prioridad, time.Now(), cas.SuperiorId,rol).Scan(&cas.Id,&cas.CreatedOn,&cas.Rol)
				if err != nil {
					log.Println(err)
				}
	}
	return
}

func (p *pgCasoRepository) fetchCasosWithClient(ctx context.Context, query string, args ...interface{}) (result []caso.Caso, err error) {
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
			&t.Rol,
			&t.ClienteName,
			&t.ClienteApellido,
			&t.FuncionarioName,
			&t.FuncionarioApellido,

		)
		result = append(result, t)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}



func (p *pgCasoRepository) fetchCasosDetail(ctx context.Context, query string, args ...interface{}) (result []caso.Caso, err error) {
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
			&t.ClienteName,
			&t.ClienteApellido,
			&t.FuncionarioName,
			&t.FuncionarioApellido,
			&t.Titulo,
			&t.Id,
			&t.Descripcion,
			&t.DetallesDeFinalizacion,
			&t.Empresa,
			&t.Area,
			&t.CreatedOn,
			&t.UpdatedOn,
			&t.FechaInicio,
			&t.FechaFin,
			&t.Prioridad,
			&t.Estado,
			&t.ClienteId,
			&t.FuncionarioId,
			&t.SuperiorId,
			&t.Rol,
			&t.ProyectoName,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}
