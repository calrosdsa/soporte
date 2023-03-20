package repository

import (
	"context"
	// "errors"
	// "soporte-go/core/model"
	"time"

	// "database/sql"
	"log"
	"soporte-go/core/model/empresa"


	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"soporte-go/core/model"
	
)

type pgEmpresaRepository struct {
	Conn    *pgxpool.Pool
	Context context.Context
}

func NewPgEmpresaRepository(conn *pgxpool.Pool, ctx context.Context) empresa.EmpresaRepository {
	return &pgEmpresaRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func (p *pgEmpresaRepository) AreaChangeState(ctx context.Context,state int,id int)(err error){
	query := `update areas set estado = $1 where id = $2;`
	_,err= p.Conn.Exec(ctx,query,state,id)
	log.Println(err)
	return 
}

func (p *pgEmpresaRepository) GetAreaByName(ctx context.Context,n string)(res empresa.Area,err error){
	query := `select * from areas where nombre = $1;`
	list,err := p.fetchAreas(ctx,query,n)
	if err != nil {
		return empresa.Area{},err
	}
	if len(list) > 0 {
		res = list[0]
	} else {
		return res, model.ErrNotFound
	}

	return

}

func(p *pgEmpresaRepository)AddUserToArea(ctx context.Context,id *string,n *string,a *empresa.AddUserRequestData) (err error){
	// query := `UPDATE clientes SET areas = areas || '{ $1 }' WHERE client_id = $2;`
	query:= `insert into user_area (user_id,area_id,nombre_user,nombre_area) values($1,$2,$3,$4);`
	_,err = p.Conn.Exec(ctx,query,id,a.AreaId,n,a.AreaName)
	return 
}

func (p *pgEmpresaRepository) GetAreasUser(ctx context.Context,userId string) (res []empresa.AreaUser, err error) {
	// var superiorId string
	// // log.Println(userId)
	// query := `select superior_id from clientes where client_id = $1;`
	// err = p.Conn.QueryRow(ctx,query,userId).Scan(&superiorId)
	// log.Println(superiorId)
	query1 := `select area_id,nombre_area,estado from user_area where user_id = $1;`
	res,err = p.fetchAreasUser(ctx,query1,userId)
	// log.Println(res)
	if err != nil {
		return nil,err
	}
	return 
}


func (p *pgEmpresaRepository) GetAreasUserAdmin(ctx context.Context,userId *string,rol *int) (res []empresa.Area, err error) {
	var superiorId string
	// log.Println(userId)
	if *rol == int(model.RoleCliente) || *rol == int(model.RoleClienteAdmin) {
		query := `select superior_id from clientes where client_id = $1;`
		err = p.Conn.QueryRow(ctx,query,userId).Scan(&superiorId)
		if err != nil {
			return nil,err
	}
	//Listar las áreas cuyo campo de 'creador_id' coincide con el campo 'superior_id' del cliente o funcionario correspondiente
	//Solo traer areas cuyo estado sean 0 (0 = disponibles o activas)
	query1 := `select * from areas where creador_id = $1 and estado = 0;`
	res,err = p.fetchAreas(ctx,query1,superiorId)
	if err != nil {
		return nil,err
	}
   }else if *rol == int(model.RoleFuncionarioAdmin) || *rol == int(model.RoleFuncionario){
	query := `select superior_id from funcionarios where funcionario_id = $1;`
	err = p.Conn.QueryRow(ctx,query,userId).Scan(&superiorId)
	if err != nil {
		return nil,err
	}
// log.Println(superiorId)
	query1 := `select * from areas where creador_id = $1 and estado = 0;`
	res,err = p.fetchAreas(ctx,query1,superiorId)
	// log.Println(res)
	if err != nil {
		return nil,err
	}
   }
	return 
}

func (p *pgEmpresaRepository) StoreEmpresa(ctx context.Context, empresa *empresa.Empresa) (err error) {
	query := `insert into empresas (nombre,slug,telefono,created_on) values($1,$2,$3,$4);`
	_,err = p.Conn.Exec(ctx,query,empresa.Nombre,empresa.Slug,empresa.Telefono,time.Now())
	return err
}

func (p *pgEmpresaRepository) GetAreasEmpresa(ctx context.Context,id int)(res []empresa.Area,err error){
	query := `select * from areas where empresa_id = $1`
	list,err := p.fetchAreas(ctx,query,id)
	if err != nil {
		return list,err
	}
	return list,err
}

func (p *pgEmpresaRepository) StoreArea(ctx context.Context, area *empresa.Area,rol *int) (id int,err error) {
	//Listar todas las áreas que pertenecen a una entidad en particular utilizando su identificador único (ID).
	list,err := p.GetAreasEmpresa(ctx,area.EmpresaId)
	if err != nil {
		log.Println(err)
		return
	}
	//Recorre la lista de áreas previamente obtenida y comprueba que el nombre de la nueva área que se quiere
	// crear no entre en conflicto con el nombre de otra área existente en la misma empresa o entidad.
	for _,item := range list{
		if item.Nombre == area.Nombre{
			return 0,model.ErrConflict
		}
	}
	var superiorId string
	//Obtener el valor del campo "superior_id" basándose en el ID del cliente o del funcionario, con el propósito
	// de utilizar este valor para crear el área correspondiente posteriormente.
	if *rol == int(model.RoleClienteAdmin) {
		log.Println(area.CreadorId)
		queryClient := `select superior_id from clientes where client_id = $1;`
		err= p.Conn.QueryRow(ctx,queryClient,area.CreadorId).Scan(&superiorId)
		if err != nil {
			return 
		}
		log.Println(superiorId)
	}else if *rol == int(model.RoleFuncionarioAdmin) {
		queryClient := `select superior_id from funcionarios where funcionario_id = $1;`
		err= p.Conn.QueryRow(ctx,queryClient,area.CreadorId).Scan(&superiorId)
		if err != nil {
			return 
		}
	}
	var areaId int
	query1 := `insert into areas (nombre,empresa_id,created_on,creador_id)
	values ($1,$2,$3,$4) returning (id);`
	err = p.Conn.QueryRow(ctx,query1,area.Nombre,area.EmpresaId	,time.Now(),superiorId).Scan(&areaId)
	if err != nil {
		return 0,err
	}

	return areaId,err
}

func (p *pgEmpresaRepository) GetEmpresa(ctx context.Context,userId string,rol int) (res empresa.Empresa, err error) {
	var empresaId int
	log.Println(userId)
	if rol == int(model.RoleFuncionario) || rol == int(model.RoleFuncionarioAdmin){
		query := `select empresa_id from funcionarios where user_id = $1;`
		err = p.Conn.QueryRow(ctx, query, userId).Scan(&empresaId)
	}else{
		query := `select empresa_id from clientes where user_id = $1;`
		err = p.Conn.QueryRow(ctx, query, userId).Scan(&empresaId)
	}
	if err != nil {
		log.Println(err)
	}
	log.Println(empresaId)
	query2 := `select * from empresas where id = $1;`
	list, err := p.fetchEmpresa(ctx, query2,empresaId)
	if err != nil {
		return empresa.Empresa{}, err
	}
	if len(list) > 0 {
		res = list[0]
	} else {
		return res, model.ErrNotFound
	}
	return
	
}

func (p *pgEmpresaRepository) GetEmpresas(ctx context.Context) (res []empresa.Empresa, err error) {
	query := `select * from empresas;`
	res,err = p.fetchEmpresa(ctx,query)
	if err != nil {
		log.Println(err)
		return nil,err
	}
	return 
}


func (p *pgEmpresaRepository) fetchEmpresa(ctx context.Context, query string, args ...interface{}) (result []empresa.Empresa, err error) {
	rows,err:= p.Conn.Query(p.Context, query, args...)
	defer func ()  {
		rows.Close()
	   }()
	result = make([]empresa.Empresa, 0)
	for rows.Next(){
		t := empresa.Empresa{}
		err = rows.Scan(
			&t.Id,
			&t.Slug,
			&t.Nombre,
			&t.Telefono,
			&t.Estado,
			&t.CreatedOn,
			&t.UpdatedOn,
		)		
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	result = append(result, t)
	log.Println(result)
}

	return result, nil
}



func (p *pgEmpresaRepository) fetchAreas(ctx context.Context, query string, args ...interface{}) (result []empresa.Area, err error) {
	rows,err:= p.Conn.Query(p.Context, query, args...)
	defer func ()  {
	 rows.Close()
	}()
	result = make([]empresa.Area, 0)
	for rows.Next(){
		t := empresa.Area{}
		err = rows.Scan(
			&t.Id,
			&t.Nombre,
			&t.Estado,
			&t.EmpresaId,
			&t.CreatedOn,
			&t.CreadorId,
		)		
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	result = append(result, t)
	log.Println(result)
}

	return result, nil
}


func (p *pgEmpresaRepository) fetchAreasUser(ctx context.Context, query string, args ...interface{}) (result []empresa.AreaUser, err error) {
	rows,err:= p.Conn.Query(p.Context, query, args...)
	defer func ()  {
	 rows.Close()
	}()
	result = make([]empresa.AreaUser, 0)
	for rows.Next(){
		t := empresa.AreaUser{}
		err = rows.Scan(
			&t.Id,
			&t.Name,
			&t.Estado,
		)		
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	result = append(result, t)
	log.Println(result)
}

	return result, nil
}
