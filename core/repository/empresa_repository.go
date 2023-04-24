package repository

import (
	"context"
	"database/sql"

	// "errors"
	// "soporte-go/core/model"
	"time"

	// "database/sql"
	"log"
	"soporte-go/core/model/empresa"

	"soporte-go/core/model"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type pgEmpresaRepository struct {
	Conn    *sql.DB
	Context context.Context
}

func NewPgEmpresaRepository(conn *sql.DB, ctx context.Context) empresa.EmpresaRepository {
	return &pgEmpresaRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func (p *pgEmpresaRepository) AreaChangeState(ctx context.Context, state int, id int) (err error) {
	query := `update areas set estado = $1 where id = $2;`
	_, err = p.Conn.ExecContext(ctx, query, state, id)
	log.Println(err)
	return
}

func (p *pgEmpresaRepository) GetAreaByName(ctx context.Context, n string) (res empresa.Area, err error) {
	var query string
	query = `select * from areas where nombre = $1;`
	list, err := p.fetchAreas(ctx, query, n)
	if err != nil {
		return empresa.Area{}, err
	}
	if len(list) > 0 {
		res = list[0]
	} else {
		return res, model.ErrNotFound
	}

	return

}

func (p *pgEmpresaRepository) GetProyectoByName(ctx context.Context, n string) (res empresa.ProyectoDetail, err error) {
	var query string
	query = `select  p.id,p.nombre,p.parent_id,p.estado ,p.empresa_id,p.empresa_parent_id ,p.created_on,p.creador_id,
     e.nombre,a.nombre from proyectos as p 
	 left join empresas as e on e.id = empresa_id
	 left join areas as a on a.id = p.parent_id
	 where p.nombre = $1; `
	err = p.Conn.QueryRowContext(ctx, query, n).Scan(&res.Id, &res.Nombre, &res.ParentId, &res.Estado, &res.EmpresaId, &res.EmpresaParentId,
		&res.CreatedOn, &res.CreadorId, &res.EmpresaName, &res.AreaName)
	if err != nil {
		log.Println(err)
		return
	}
	query = `select id,proyecto_id,start_date,end_date from proyecto_duration where proyecto_id = $1`
	list, err := p.fetchProyectoDuration(ctx, query, res.Id)

	res.ProyectoDuration = list
	return

}

func (p *pgEmpresaRepository) AddUserToArea(ctx context.Context, id string, a *empresa.AddUserRequestData) (err error) {
	// query := `UPDATE clientes SET areas = areas || '{ $1 }' WHERE client_id = $2;`
	query := `insert into user_area (user_id,area_id,nombre_area,created_on) values($1,$2,$3,$4);`
	_, err = p.Conn.ExecContext(ctx, query, id, a.AreaId, a.AreaName, time.Now())
	return
}

// Lista todos los user_area
func (p *pgEmpresaRepository) GetAreasUser(ctx context.Context, userId string) (res []empresa.AreaUser, err error) {

	query1 := `select area_id,nombre_area,estado from user_area where user_id = $1;`
	res, err = p.fetchAreasUser(ctx, query1, userId)
	if err != nil {
		return nil, err
	}
	return
}

func (p *pgEmpresaRepository) GetClientesByAreaId(ctx context.Context, areaId int) (res []empresa.UserArea, err error) {
	query := `select client_id,nombre,apellido,profile_photo,user_area.estado from clientes inner join user_area 
	on clientes.client_id = user_area.user_id where user_area.area_id = $1;`
	res, err = p.fetchUserArea(ctx, query, areaId)
	if err != nil {
		return
	}
	return
}

func (p *pgEmpresaRepository) GetFuncionariosByAreaId(ctx context.Context, areaId int) (res []empresa.UserArea, err error) {
	query := `select funcionario_id,nombre,apellido,profile_photo,user_area.estado from funcionarios inner join user_area 
	on funcionarios.funcionario_id = user_area.user_id where user_area.area_id = $1;`
	res, err = p.fetchUserArea(ctx, query, areaId)
	if err != nil {
		return
	}
	return
}

func (p *pgEmpresaRepository) GetProyectoFromUserArea(ctx context.Context, id string) (res []empresa.Area, err error) {
	query := `select area_id,nombre_area,estado,(0),created_on,user_id from user_area where user_id = $1`
	res, err = p.fetchAreas(ctx, query, id)
	if err != nil {
		return
	}
	return
}

func (p *pgEmpresaRepository) GetProyectosFuncionario(ctx context.Context, id string) (res []empresa.Area, err error) {
	query := `select id,nombre,estado,empresa_id,created_on,creador_id from proyectos where creador_id = $1`
	res, err = p.fetchAreas(ctx, query, id)
	if err != nil {
		return
	}
	return
}

func (p *pgEmpresaRepository) GetProyectosAdmin(ctx context.Context, emId int) (res []empresa.Area, err error) {
	query := `select id,nombre,estado,empresa_id,created_on,creador_id from proyectos where empresa_parent_id = $1`
	res, err = p.fetchAreas(ctx, query, emId)
	if err != nil {
		return
	}
	return
}

func (p *pgEmpresaRepository) StoreEmpresa(ctx context.Context, empresa *empresa.Empresa) (err error) {
	query := `insert into empresas (nombre,slug,telefono,created_on,parent_id) values($1,$2,$3,$4,$5);`
	_, err = p.Conn.ExecContext(ctx, query, empresa.Nombre, empresa.Slug, empresa.Telefono, time.Now(), empresa.ParentId)
	return err
}

func (p *pgEmpresaRepository) GetAreasEmpresa(ctx context.Context, id int) (res []empresa.Area, err error) {
	query := `select * from areas where empresa_id = $1`
	res, err = p.fetchAreas(ctx, query, id)
	if err != nil {
		return
	}
	return
}

func (p *pgEmpresaRepository) GetAreasFuncionario(ctx context.Context, id string) (res []empresa.Area, err error) {
	query := `select * from areas where creador_id = $1`
	res, err = p.fetchAreas(ctx, query, id)
	if err != nil {
		return
	}
	return
}

func (p *pgEmpresaRepository) GetProyectosEmpresa(ctx context.Context, emId int) (res []empresa.Area, err error) {
	query := `select id,nombre,estado,empresa_id,created_on,creador_id from proyectos where empresa_id = $1`
	res, err = p.fetchAreas(ctx, query, emId)
	if err != nil {
		return
	}
	return
}

func (p *pgEmpresaRepository) GetProyectos(ctx context.Context, parentId int) (res []empresa.Area, err error) {
	query := `select id,nombre,estado,empresa_id,created_on,creador_id from proyectos where parent_id = $1`
	// log.Println(parentIdem)
	res, err = p.fetchAreas(ctx, query, parentId)
	// if err != nil {
	// return
	// }
	return
}

func (p *pgEmpresaRepository) CreateProyecto(ctx context.Context, a *empresa.Proyecto) (err error) {
	var query string
	list, err := p.GetProyectosEmpresa(ctx, a.EmpresaId)
	if err != nil {
		return
	}
	for _, item := range list {
		if item.Nombre == a.Nombre {
			return model.ErrConflict
		}
	}
	var proyectoId int
	query = `insert into proyectos (nombre,empresa_id,created_on,creador_id,empresa_parent_id,parent_id)
	values ($1,$2,$3,$4,$5,$6) returning (id)`
	err = p.Conn.QueryRowContext(ctx, query, a.Nombre, a.EmpresaId, time.Now(), a.CreadorId, a.EmpresaParentId, a.ParentId).Scan(&proyectoId)
	if err != nil {
		return
	}

	start, _ := time.Parse("2006-01-02", a.Start)
	end, _ := time.Parse("2006-01-02", a.End)
	query = `insert into proyecto_duration (proyecto_id,start_date,end_date) values($1,$2,$3)`
	_, err = p.Conn.ExecContext(ctx, query, proyectoId, start, end)
	if err != nil {
		log.Println(err)
	}

	// log.Println(res.LastInsertId())
	// a.Id,_ = res.LastInsertId()
	query = `insert into user_area (user_id,area_id,nombre_area,created_on) values($1,$2,$3,$4);`
	_, err = p.Conn.ExecContext(ctx, query, a.CreadorId, proyectoId, a.Nombre, time.Now())
	if err != nil {
		return
	}
	a.Id = proyectoId
	return
}

func (p *pgEmpresaRepository) StoreArea(ctx context.Context, area *empresa.Area) (err error) {
	//Listar todas las áreas que pertenecen a una entidad en particular utilizando su identificador único (ID).
	list, err := p.GetAreasEmpresa(ctx, area.EmpresaId)
	if err != nil {
		// log.Println(err)
		return
	}
	//Recorre la lista de áreas previamente obtenida y comprueba que el nombre de la nueva área que se quiere
	// crear no entre en conflicto con el nombre de otra área existente en la misma empresa o entidad.
	for _, item := range list {
		if item.Nombre == area.Nombre {
			return model.ErrConflict
		}
	}
	var areaId int
	query1 := `insert into areas (nombre,empresa_id,created_on,creador_id)
	values ($1,$2,$3,$4) returning (id)`
	err = p.Conn.QueryRowContext(ctx, query1, area.Nombre, area.EmpresaId, time.Now(), area.CreadorId).Scan(&areaId)
	// log.Println(res.LastInsertId())
	// area.Id,_ = res.LastInsertId()
	if err != nil {
		return
	}
	area.Id = areaId
	return
}

func (p *pgEmpresaRepository) GetEmpresa(ctx context.Context, userId string, rol int) (res empresa.Empresa, err error) {
	var empresaId int
	log.Println(userId)
	if rol == int(model.RoleFuncionario) || rol == int(model.RoleFuncionarioAdmin) {
		query := `select empresa_id from funcionarios where user_id = $1;`
		err = p.Conn.QueryRowContext(ctx, query, userId).Scan(&empresaId)
	} else {
		query := `select empresa_id from clientes where user_id = $1;`
		err = p.Conn.QueryRowContext(ctx, query, userId).Scan(&empresaId)
	}
	if err != nil {
		log.Println(err)
	}
	log.Println(empresaId)
	query2 := `select * from empresas where id = $1;`
	list, err := p.fetchEmpresa(ctx, query2, empresaId)
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

func (p *pgEmpresaRepository) GetEmpresas(ctx context.Context, emId *int) (res []empresa.Empresa, err error) {
	query := `select * from empresas where parent_id = $1`
	res, err = p.fetchEmpresa(ctx, query, emId)
	if err != nil {
		// log.Println(err)
		return nil, err
	}
	return
}

func (p *pgEmpresaRepository) fetchEmpresa(ctx context.Context, query string, args ...interface{}) (result []empresa.Empresa, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	defer func() {
		rows.Close()
	}()
	result = make([]empresa.Empresa, 0)
	for rows.Next() {
		t := empresa.Empresa{}
		err = rows.Scan(
			&t.Id,
			&t.Slug,
			&t.Nombre,
			&t.Telefono,
			&t.Estado,
			&t.CreatedOn,
			&t.UpdatedOn,
			&t.ParentId,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
		// log.Println(result)
	}

	return result, nil
}

func (p *pgEmpresaRepository) fetchAreas(ctx context.Context, query string, args ...interface{}) (result []empresa.Area, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	defer func() {
		rows.Close()
	}()
	result = make([]empresa.Area, 0)
	for rows.Next() {
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

func (p *pgEmpresaRepository) fetchProyectoDuration(ctx context.Context, query string, args ...interface{}) (result []empresa.ProyectoDuration, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	defer func() {
		rows.Close()
	}()
	result = make([]empresa.ProyectoDuration, 0)
	for rows.Next() {
		t := empresa.ProyectoDuration{}
		err = rows.Scan(
			&t.Id,
			&t.ProyectoId,
			&t.StartDate,
			&t.EndDate,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (p *pgEmpresaRepository) fetchAreasUser(ctx context.Context, query string, args ...interface{}) (result []empresa.AreaUser, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	defer func() {
		rows.Close()
	}()
	result = make([]empresa.AreaUser, 0)
	for rows.Next() {
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

func (p *pgEmpresaRepository) fetchUserArea(ctx context.Context, query string, args ...interface{}) (result []empresa.UserArea, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	defer func() {
		rows.Close()
	}()
	result = make([]empresa.UserArea, 0)
	for rows.Next() {
		t := empresa.UserArea{}
		err = rows.Scan(
			&t.UserId,
			&t.Name,
			&t.Apellido,
			&t.Photo,
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
