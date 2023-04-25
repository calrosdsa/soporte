package caso

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	// "database/sql"
	"log"
	// "soporte-go/core/model"
	"soporte-go/core/model"
	"soporte-go/core/model/caso"
	"soporte-go/core/model/user"
	"soporte-go/core/model/ws"

	// "github.com/jackc/pgx/v5/pgxpool"
	// "github.com/aws/aws-sdk-go/private/protocol/query"
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

func (p *pgCasoRepository) GetCasosCliForReporte(ctx context.Context, options *caso.CasoReporteOptions) (res []caso.Caso, err error) {
	var query string
	log.Println(options.Estados)
	if len(options.Estados) == 3 {
		// log.Println("ALL estaods")
		query = `select clientes.nombre,clientes.apellido,funcionarios.nombre,funcionarios.apellido,titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
		c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.client_id,c.funcionario_id,c.superior_id,c.rol,p.nombre,c.key
		from casos as c inner join clientes on clientes.client_id = c.client_id left join funcionarios
		 on funcionarios.funcionario_id = c.funcionario_id
		left join proyectos as p on c.area = p.id
		 where c.created_on between $1 and $2 and c.area = any($3);`
		res, err = p.fetchCasosDetail(ctx, query, options.StartDate, options.EndDate, pq.Array(options.Areas))
	} else {
		query = `select clientes.nombre,clientes.apellido,funcionarios.nombre,funcionarios.apellido,titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
		c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.client_id,c.funcionario_id,c.superior_id,c.rol,p.nombre,c.key
		from casos as c inner join clientes on clientes.client_id = c.client_id 
		left join funcionarios on funcionarios.funcionario_id = c.funcionario_id
		left join proyectos as p on c.area = p.id
		 where c.created_on between $1 and $2 and c.estado = any($3) and c.area = any($4);`
		res, err = p.fetchCasosDetail(ctx, query, options.StartDate, options.EndDate, pq.Array(options.Estados), pq.Array(options.Areas))
	}
	return
}

func (p *pgCasoRepository) GetCasosFunForReporte(ctx context.Context, options *caso.CasoReporteOptions) (res []caso.Caso, err error) {
	var query string
	log.Println(options.Estados)
	if len(options.Estados) == 3 {
		// log.Println("ALL estaods")
		query = `select uc.nombre,uc.apellido,uf.nombre,uf.apellido,titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
		c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.funcionario_id,c.funcionario_id,c.superior_id,c.rol,p.nombre,c.key
		from casos as c inner join funcionarios as uc on uc.funcionario_id = c.client_id 
		left join funcionarios as uf on uf.funcionario_id = c.funcionario_id
		left join proyectos as p on c.area = p.id
		 where c.created_on between $1 and $2 and c.area = any($3);`
		res, err = p.fetchCasosDetail(ctx, query, options.StartDate, options.EndDate, pq.Array(options.Areas))
	} else {
		query = `select uc.nombre,uc.apellido,uf.nombre,uf.apellido,titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
		c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.funcionario_id,c.funcionario_id,c.superior_id,c.rol,p.nombre,c.key
		from casos as c inner join funcionarios as uc on uc.funcionario_id = c.client_id 
		left join funcionarios as uf on uf.funcionario_id = c.funcionario_id
		left join proyectos as p on c.area = p.id
		 where c.created_on between $1 and $2 and c.estado = any($3) and c.area = any($4);`
		res, err = p.fetchCasosDetail(ctx, query, options.StartDate, options.EndDate, pq.Array(options.Estados), pq.Array(options.Areas))
	}
	return
}

func (p *pgCasoRepository) FinalizarCaso(ctx context.Context, fD *caso.FinalizacionDetail) (err error) {
	log.Println(fD.Detail)
	query := `update casos set detalles_de_finalizacion = $1, estado = $2,fecha_fin = $3,updated_on = $3 where id = $4`
	_, err = p.Conn.ExecContext(ctx, query, fD.Detail, fD.Estado, time.Now(), fD.Id)
	return
}

func (p *pgCasoRepository) AsignarFuncionarioSoporte(ctx context.Context, id string, uId string) (err error) {
	// var user []user.UserForList
	query := `insert into usuarios_caso (caso_id,user_id) values ($1,$2)`
	_, err = p.Conn.ExecContext(ctx, query, id, uId)
	if err != nil {
		log.Println(err)
	}
	return

}

func (p *pgCasoRepository) AsignarFuncionario(ctx context.Context, id string, idF string) (err error) {
	var query string
	var fechaInicio *string

	query = `update casos set funcionario_id = $1,fecha_inicio = $2,updated_on = $3 where id = $4 returning (fecha_inicio)`
	err = p.Conn.QueryRowContext(ctx, query, idF, time.Now(), time.Now(), id).Scan(&fechaInicio)
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
	var query string
	query = `select clientes.nombre,clientes.apellido,funcionarios.nombre,funcionarios.apellido,titulo,casos.id,descripcion,detalles_de_finalizacion,empresa,area,casos.created_on,
	casos.updated_on,fecha_inicio,fecha_fin,prioridad,casos.estado,casos.client_id,casos.funcionario_id,casos.superior_id,
	casos.rol,p.nombre,casos.key
	from casos inner join clientes on clientes.client_id = casos.client_id 
	left join funcionarios on funcionarios.funcionario_id = casos.funcionario_id
	left join proyectos as p on casos.area = p.id
	where casos.id = $1 limit 1;`

	// var casoDetail caso.Caso
	list, err := p.fetchCasosDetail(ctx, query, id)
	if err != nil {
		return caso.Caso{}, err
	}
	res = list[0]
	query = `select * from usuarios_caso where caso_id = $1`
	// res.UsuariosCaso,err = p.fetchUsersForList(ctx,query,res.Id)
	return
}

func (p *pgCasoRepository) GetUsuariosCaso(ctx context.Context, cId string) (res []user.UserForList, err error) {
	query := `select f.funcionario_id,f.nombre,f.apellido,f.profile_photo from usuarios_caso as uc
	left join funcionarios as f on f.funcionario_id = uc.user_id  where caso_id = $1`
	res, err = p.fetchUsersForList(ctx, query, cId)
	if err != nil {
		log.Println(err)
	}
	return
}

func (p *pgCasoRepository) GetCasoFuncionario(ctx context.Context, id string) (res caso.Caso, err error) {
	query := `select 
	uc.nombre,uc.apellido,
	uf.nombre,uf.apellido,
	titulo,c.id,descripcion,detalles_de_finalizacion,empresa,area,c.created_on,
	c.updated_on,fecha_inicio,fecha_fin,prioridad,c.estado,c.client_id,c.funcionario_id,c.superior_id,c.rol,p.nombre,c.key
	from casos as c left join funcionarios as uc on uc.funcionario_id = c.client_id 
	left join proyectos as p on c.area = p.id
	left join funcionarios as uf on uf.funcionario_id = c.funcionario_id
	where c.id = $1 limit 1`
	// var casoDetail caso.Caso
	list, err := p.fetchCasosDetail(ctx, query, id)
	if err != nil {
		log.Println(err)
		return caso.Caso{}, err
	}
	res = list[0]
	// query = `select * from usuarios_caso where caso_id = $1`
	// res.UsuariosCaso,err = p.fetchUsersForList(ctx,query,res.Id)
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

func (p *pgCasoRepository) GetCasosCliente(ctx context.Context, id string, q *caso.CasoQuery) (list []caso.Caso, err error) {
	// var count int
	// query2 := `select id,titulo,created_on,updated_on,prioridad,estado,rol
	// from casos where client_id = $1 limit $2 offset $3`
	query2 := fmt.Sprintf(`select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	cl.nombre,cl.apellido,
	f2.nombre,f2.apellido,key
	from casos as c 
	left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	where c.client_id = $1 %s %s %s limit $2 offset $3`, q.Key, q.Proyecto, q.Order)
	list, _ = p.fetchCasosWithClient(ctx, query2, id, q.PageSize, q.Page*10)
	// query3 := `select count(*) from casos where client_id = $1;`
	// err = p.Conn.QueryRowContext(ctx,query3,id).Scan(&count)
	// size = count
	return
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetCasosFuncionario(ctx context.Context, id string, q *caso.CasoQuery) (res []caso.Caso, err error) {
	// query := fmt.Sprintf(`
	// select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	// coalesce(cl.nombre,f.nombre) as nombre,
	// coalesce(cl.apellido,f.apellido) as apellido,
	// f2.nombre,f2.apellido,key
	// from casos as c 
	// left join clientes as cl on cl.client_id = c.client_id 
	// left join funcionarios as f on f.funcionario_id = c.client_id
	// left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	// where c.funcionario_id = $1 or c.client_id = $1 %s %s 
	// union 
	// select c.id,c.titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	// coalesce(cl.nombre,f.nombre) as nombre,
	// coalesce(cl.apellido,f.apellido) as apellido,
	// f2.nombre,f2.apellido,key
	// from usuarios_caso as uc inner join casos as c on c.id = uc.caso_id
	// left join clientes as cl on cl.client_id = c.client_id 
	// left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	// left join funcionarios as f on f.funcionario_id = c.client_id
	// where  uc.user_id = $1 %s %s
	// %s limit $2 offset $3`, q.Key, q.Proyecto,q.Key,q.Proyecto, q.Order)
	// res, err = p.fetchCasosWithClient(ctx, query, id, q.PageSize, q.Page*10)
	// if err != nil {
	// 	log.Println(err)
	// }
	// return
	query1 := fmt.Sprintf(`
	select c.id,c.titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido,
	f2.nombre,f2.apellido,key
	from usuarios_caso as uc inner join casos as c on c.id = uc.caso_id
	left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	left join funcionarios as f on f.funcionario_id = c.client_id
	where  uc.user_id = $1 %s %s
	union
	select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido,
	f2.nombre,f2.apellido,key
	from casos as c 
	left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f on f.funcionario_id = c.client_id
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	where c.funcionario_id = $1 %s %s 
	union
	select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido,
	f2.nombre,f2.apellido,key
	from casos as c 
	left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f on f.funcionario_id = c.client_id
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	where c.client_id = $1 %s %s 
	%s limit $2 offset $3`,q.Key, q.Proyecto, q.Key, q.Proyecto,q.Key, q.Proyecto, q.Order)
	res, err = p.fetchCasosWithClient(ctx, query1,id, q.PageSize, q.Page*10)
	log.Println(err)
	return
}

func (p *pgCasoRepository) GetAllCasosUserCliente(ctx context.Context, id string, q *caso.CasoQuery) (list []caso.Caso, err error) {
	// var superiorId string
	// query1 := `select superior_id from clientes where client_id = $1;`
	// if err = p.Conn.QueryRowContext(ctx, query1, id).Scan(&superiorId); err != nil {
	// 	return
	// }
	query2 := fmt.Sprintf(`select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	cl.nombre,cl.apellido,
	f2.nombre,f2.apellido,key
	from casos as c left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	where c.superior_id = $1 %s %s %s limit 10 offset $2`, q.Key, q.Proyecto, q.Order)
	if q.Page == 1 || q.Page == 0 {
		list, err = p.fetchCasosWithClient(ctx, query2, id, 0)
	} else {
		page := q.Page - 1
		list, err = p.fetchCasosWithClient(ctx, query2, id, page*10)
	}
	log.Println(err)

	return
	// return []caso.Caso{},nil
}

func (p *pgCasoRepository) GetAllCasosUserFuncionario(ctx context.Context, id int, q *caso.CasoQuery) (res []caso.Caso, err error) {
	// log.Println("Get all casos funcionario")
	// log.Println(model.ASC)
	// pq.
	query2 := fmt.Sprintf(`select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido,
	f2.nombre,f2.apellido,key
	from casos as c left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	left join funcionarios as f on f.funcionario_id = c.client_id where status = 0 %s %s
	%s limit $1 offset $2`, q.Key, q.Proyecto, q.Order)
	res, err = p.fetchCasosWithClient(ctx, query2, q.PageSize, q.Page*10)
	return
}

func (p *pgCasoRepository) GetCasosFromUserCaso(ctx context.Context, id string, q *caso.CasoQuery) (res []caso.Caso, err error) {
	log.Println(q.Key)
	log.Println(q.Proyecto)
	log.Println(q.Order)

	query := fmt.Sprintf(`select c.id,c.titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido,
	f2.nombre,f2.apellido,key
	from usuarios_caso as uc inner join casos as c on c.id = uc.caso_id
	left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	left join funcionarios as f on f.funcionario_id = c.client_id
	where  uc.user_id = $1 and status = 0 %s %s
	union 
	select id,titulo,c.created_on,c.updated_on,prioridad,c.estado,c.rol,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido,
	f2.nombre,f2.apellido,key
	from casos as c 
	left join clientes as cl on cl.client_id = c.client_id 
	left join funcionarios as f on f.funcionario_id = c.client_id
	left join funcionarios as f2 on f2.funcionario_id = c.funcionario_id
	where c.funcionario_id = $1 or c.client_id = $1
	%s %s %s limit $2 offset $3`, q.Key, q.Proyecto,q.Key,q.Proyecto, q.Order)
	res, err = p.fetchCasosWithClient(ctx, query, id, q.PageSize, q.Page*10)
	if err != nil {
		log.Println(err)
	}
	return
}

func (p *pgCasoRepository) UpdateCaso(ctx context.Context, c *caso.Caso) (err error) {
	query := `update casos set titulo = $1,descripcion = $2,updated_on = $3 where id = $4 
	returning titulo,descripcion,updated_on`
	err = p.Conn.QueryRowContext(ctx, query, c.Titulo, c.Descripcion, time.Now(), c.Id).Scan(
		&c.Titulo, &c.Descripcion, &c.UpdatedOn)
	if err != nil {
		log.Println(err)
	}
	return
}

func (p *pgCasoRepository) CreateCasoCliente(ctx context.Context, cas *caso.Caso, id string, emI int, rol int) (err error) {
	rand.Seed(time.Now().UnixNano())
	num := strconv.Itoa(rand.Intn(9000) + 1000)
	key := fmt.Sprintf("%s-%s", cas.Key, num)
	query := `select superior_id from clientes where client_id = $1;`
	err = p.Conn.QueryRowContext(ctx, query, id).Scan(&cas.SuperiorId)
	if err != nil {
		return
	}
	query1 := `insert into casos (titulo,client_id,descripcion,empresa,area,prioridad,created_on,superior_id,rol,key) values(
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id,created_on,rol,key,estado;`
	err = p.Conn.QueryRowContext(ctx, query1, cas.Titulo, id, cas.Descripcion, emI, cas.Area,
		cas.Prioridad, time.Now(), cas.SuperiorId, rol, key).Scan(&cas.Id, &cas.CreatedOn, &cas.Rol, &cas.Key, &cas.Estado)
	if err != nil {
		log.Println(err)
	}
	return
}

func (p *pgCasoRepository) CreateCasoFuncionario(ctx context.Context, cas *caso.Caso, id string, emI int, rol int) (err error) {
	log.Println("Create casp for funcionario")
	rand.Seed(time.Now().UnixNano())
	num := strconv.Itoa(rand.Intn(9000) + 1000)
	key := fmt.Sprintf("%s-%s", cas.Key, num)
	query := `select superior_id from funcionarios where funcionario_id = $1;`
	err = p.Conn.QueryRowContext(ctx, query, id).Scan(&cas.SuperiorId)
	if err != nil {
		return
	}
	if rol == int(model.RoleFuncionario) {
		query1 := `insert into casos (titulo,client_id,descripcion,empresa,area,prioridad,superior_id,
			funcionario_id,rol,key,created_on,fecha_inicio) values(
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$11) returning id,created_on,rol,key,estado;`
		err = p.Conn.QueryRowContext(ctx, query1, cas.Titulo, id, cas.Descripcion, emI, cas.Area,
			cas.Prioridad, cas.SuperiorId, cas.FuncionarioId, rol, key,time.Now()).Scan(&cas.Id, &cas.CreatedOn, &cas.Rol,
			&cas.Key, &cas.Estado)
		if err != nil {
			log.Println(err)
		}
	} else {
		query1 := `insert into casos (titulo,client_id,descripcion,empresa,area,prioridad,superior_id,rol,key,
			fecha_inicio,created_on,funcionario_id) values(
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10,$11) returning id,created_on,rol,key,estado;`
		err = p.Conn.QueryRowContext(ctx, query1, cas.Titulo, id, cas.Descripcion, emI, cas.Area,
			cas.Prioridad, cas.SuperiorId, rol, key,time.Now(),cas.FuncionarioId).Scan(&cas.Id, &cas.CreatedOn, &cas.Rol, &cas.Key, &cas.Estado)
		if err != nil {
			log.Println(err)
		}
	}
	return
}


func (p *pgCasoRepository) GetMessagesCaso(ctx context.Context, casoId string) (res []ws.Message, err error) {
	query := `select m.id,m.caso_id,m.from_user,m.media_url,m.content,m.is_read,m.created_on,m.is_deleted,
	coalesce(cl.nombre,f.nombre) as nombre,
	coalesce(cl.apellido,f.apellido) as apellido
	from messages as m left join funcionarios as f on f.funcionario_id = from_user
	left join clientes as cl on cl.client_id = from_user
	where caso_id = $1;`
	res, err = p.fetchMessages(ctx, query, casoId)

	return
}
func (p *pgCasoRepository) fetchMessages(ctx context.Context, query string, args ...interface{}) (result []ws.Message, err error) {
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
			pq.Array(&t.MediaUrl),
			&t.Content,
			&t.IsRead,
			&t.CreatedOn,
			&t.IsDeleted,
			&t.Nombre,
			&t.Apellido,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
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
			&t.Key,
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
			&t.Key,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}

func (p *pgCasoRepository) fetchUsersForList(ctx context.Context, query string, args ...interface{}) (result []user.UserForList, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]user.UserForList, 0)
	for rows.Next() {
		t := user.UserForList{}
		err = rows.Scan(
			&t.Id,
			&t.Nombre,
			&t.Apellido,
			&t.Photo,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}
