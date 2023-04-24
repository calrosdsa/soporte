package repository

import (
	"context"
	"database/sql"
	"time"

	// "databasx.Conn"
	"fmt"
	"log"
	"soporte-go/core/model"
	user "soporte-go/core/model/user"
	"soporte-go/util"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type pgUserRepository struct {
	Conn    *sql.DB
	Context context.Context
}

func NewPgUserRepository(conn *sql.DB, ctx context.Context) user.UserRepository {
	return &pgUserRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func(p *pgUserRepository) GetFuncionariosEmpresa(ctx context.Context,emId int)(res []user.UserForList,err error){
	// log.Println(emId)	
	query := `select funcionario_id,nombre,apellido,profile_photo
	from funcionarios where empresa_id = $1 and estado = 0;`
	res, err = p.fetchUsersForList(ctx, query, emId)
	return

}


func (p *pgUserRepository) GetClientesEmpresa(ctx context.Context,emId int)(res []user.UserForList,err error){
	query := `select client_id,nombre, apellido,profile_photo
	from clientes where empresa_id = $1 and estado = 0`
	res, err = p.fetchUsersForList(ctx, query, emId)
	return
}

func (p *pgUserRepository) SearchUser(ctx context.Context, id string, q string) (res []user.UserShortInfo, err error) {
	log.Println(q)
	query := `select client_id,nombre,apellido,is_admin,(false),email,profile_photo,estado,created_on
	from clientes where superior_id = $1 and nombre ILIKE $2;`
	res, err = p.fetchUserShortInfo(ctx, query, id, q)
	return

}

func (p *pgUserRepository) DeleteInvitation(ctx context.Context, m string) (err error) {
	log.Println(m)
	query := `delete from invitaciones where email = $1;`
	_, err = p.Conn.ExecContext(ctx, query, m)
	return
}

func (p *pgUserRepository) ValidateEmail(ctx context.Context, m string) (err error) {
	var email string
	query := `select email from users where email = $1;`
	p.Conn.QueryRowContext(ctx, query, m).Scan(&email);
	if email == m {
		return fmt.Errorf("ya existe un usuario con este email: %s", email)
	}
	query1 := `select email from invitaciones where email = $1;`
	p.Conn.QueryRowContext(ctx, query1, m).Scan(&email);

	if email == m {
		return fmt.Errorf("ya existe una invitacion con este email: %s", email)
	}
	return nil
}

func (p *pgUserRepository) CreateUserInvitationC(ctx context.Context, us *user.UserShortInfo) (res user.UserShortInfo, err error) {
	// var superiorId *string
	// query := `select superior_id from clientes where client_id = $1`
	// log.Println(us.Id)
	// p.Conn.QueryRowContext(ctx, query, us.Id).Scan(&superiorId)
	// if err != nil {
	// 	log.Println("error here1")
	// 	return user.UserShortInfo{}, err
	// }
	query2 := `insert into invitaciones (email,is_admin,creador_id,send_on) values($1,$2,$3,$4)
	 returning id,email,pendiente,is_admin,send_on;`
	t := user.UserShortInfo{}
	err = p.Conn.QueryRowContext(ctx, query2, us.Nombre, us.IsAdmin, us.Id,time.Now()).Scan(
		&t.Id,&t.Nombre,&t.Pendiente,&t.IsAdmin,&t.DateTime)
	if err != nil {
		log.Println("error is here",err)
		return user.UserShortInfo{}, err
	}
	return t, nil
}

func (p *pgUserRepository) CreateUserInvitationF(ctx context.Context, us *user.UserShortInfo) (res user.UserShortInfo, err error) {
	// var superiorId string
	// query := `select superior_id from funcionarios where funcionario_id = $1`
	// err = p.Conn.QueryRowContext(ctx, query, us.Id).Scan(&superiorId)
	// if err != nil {
	// 	return user.UserShortInfo{}, err
	// }
	query2 := `insert into invitaciones (email,is_admin,creador_id,send_on) values($1,$2,$3,$4)
	 returning id,email,pendiente,is_admin,send_on;;`
	t := user.UserShortInfo{}
	err = p.Conn.QueryRowContext(ctx, query2, us.Nombre, us.IsAdmin, us.Id,time.Now()).Scan(
		&t.Id,&t.Nombre,&t.Pendiente,&t.IsAdmin,&t.DateTime)
	if err != nil {
		log.Println("error is here",err)
		return user.UserShortInfo{}, err
	}
	return t, nil
}


func (p *pgUserRepository) GetClientesEmpresaByRol(ctx context.Context,emId int,rol int)(res []user.UserShortInfo,err error){
	query := `select client_id,nombre, apellido,is_admin,(false),email,profile_photo,estado,created_on
	from clientes where empresa_id = $1 and rol = $2`
	res, err = p.fetchUserShortInfo(ctx, query, emId,rol)
	return
}


func (p *pgUserRepository) GetUsersShortIInfoC(ctx context.Context, id string) (res []user.UserShortInfo, err error) {
		query := `select client_id,nombre, apellido,is_admin,(false),email,profile_photo,estado,created_on
		 from clientes where superior_id = $1;`
		res, err = p.fetchUserShortInfo(ctx, query, id)
		if err != nil {
			return
		}
	return
}

func (p *pgUserRepository) GetUsersShortIInfoF(ctx context.Context,emId int)(res []user.UserShortInfo,err error){
	query := `select funcionario_id,nombre, apellido,(false),(false),email,profile_photo,estado,created_on
		 from funcionarios where empresa_id = $1`
	res, err = p.fetchUserShortInfo(ctx, query, emId)
	if err != nil {
		return 
	}
	return
}

func (p *pgUserRepository) GetInvitaciones(ctx context.Context, id string) (res []user.UserShortInfo, err error) {
	query := `select id,email,(''),is_admin,pendiente,(''),(''),(0),send_on from invitaciones where creador_id = $1`
	res, err = p.fetchUserShortInfo(ctx, query, id)
	if err != nil {
		return nil, err
	}
	return
}



func filter[T any](ss []T, test func(T) bool) (ret []T) {
    for _, s := range ss {
        if test(s) {
            ret = append(ret, s)
        }
    }
    return
}

func (p *pgUserRepository) GetUserAddList(ctx context.Context, f int,rol int,sId string) (res []user.UserForList, err error) {
	// log.Println(f)
	var query string
	query = `select user_id,nombre_user,estado from user_area where area_id = $1;`
	res2, _ := p.fetchUserArea(ctx, query, f)
	log.Println("len res2", len(res2))
	if rol == int(model.RoleClienteAdmin){
		query = `select client_id,nombre,estado from clientes where superior_id = $1 ;`
		res, _ = p.fetchUserArea(ctx, query,sId)
	} else if rol == int(model.RoleFuncionarioAdmin){
		query = `select funcionario_id,nombre,estado from funcionarios;`
		res, _ = p.fetchUserArea(ctx, query)
	}
	if len(res2) == 0 {
		return res, nil
	} else {
		users:= res
		for i := len(res) -1;i>= 0;i-- {
		for _, val2 := range res2 {
				if val2.Id == res[i].Id {
					users[i] = users[len(users)-1] // Copy last element to index i.
					users[len(users)-1] = user.UserForList{}   // Erase last element (write zero value).
					users = users[:len(users)-1]   // Truncate slice.
				}
			}

		}
		return res, nil
	}
}

func (p *pgUserRepository) GetClientes(ctx context.Context) (clietes []user.Cliente, err error) {
	query := `select * from clientes;`
	log.Println(time.Now())
	list, err := p.fetchClientes(ctx, query)
	if err != nil {
		log.Println(err)
	}
	return list, err
}

func (p *pgUserRepository) UpdateCliente(ctx context.Context, columns []string, values ...interface{}) error {
	query, _ := util.AppendQueries(`client_id`, `update clientes set `, columns)
	
	log.Println(query)
	log.Println(values...)
	_, err := p.Conn.ExecContext(p.Context, query, values...)
	if err != nil {
		return err
	}

	// query := `update clientes set nombre = &1,apellido = &2,email=&3,celular=&4,telefono=&5,profile_photo=&6)`
	return nil
}

func (p *pgUserRepository) UpdateFuncionario(ctx context.Context, columns []string, values ...interface{}) error {
	query, _ := util.AppendQueries(`funcionario_id`, `update funcionarios set `, columns)
	log.Println(query)
	_, err := p.Conn.ExecContext(p.Context, query, values...)
	if err != nil {
		return err
	}

	// query := `update clientes set nombre = &1,apellido = &2,email=&3,celular=&4,telefono=&5,profile_photo=&6)`
	return nil
}
func (p *pgUserRepository) GetFuncionarios(ctx context.Context) (funcionarios []user.Funcionario, err error) {
	query := `select * from funcionarios;`
	list, err := p.fetchFuncionarios(ctx, query)
	if err != nil {
		log.Println(err)
	}
	return list, err
}

func (p *pgUserRepository) GetClienteDetail(ctx context.Context, id string) (res user.UserDetail, err error) {
		query := `select client_id,nombre,apellido,celular,email,superior_id,empresa_id,telefono,created_on,
		updated_on,user_id,estado,profile_photo,rol
		 from clientes where client_id = $1 limit 1;`
		log.Println(id)
		list, err := p.fetchUserDetail(ctx, query, id)
		if err != nil {
			return 
		}
		if len(list) > 0 {
			res = list[0]
		} else {
			return res, model.ErrNotFound
		}
	return
}

func (p *pgUserRepository) GetFuncionarioDetail(ctx context.Context, id string) (res user.UserDetail, err error) {
	query := `select funcionario_id,nombre,apellido,celular,email,superior_id,empresa_id,telefono,created_on,
	updated_on,user_id,estado,profile_photo,rol from funcionarios where funcionario_id = $1 limit 1;`
	// log.Println(id)
	list, err := p.fetchUserDetail(ctx, query, id)
	if err != nil {
		return
	}
	if len(list) > 0 {
		res = list[0]
	} else {
		return 
	}
return
}

func (p *pgUserRepository) GetFuncionarioById(ctx context.Context, id string) (res user.Funcionario, err error) {
	query := `select * from clientes where client_id = $1;`
	list, err := p.fetchFuncionarios(ctx, query, id)
	if err != nil {
		return user.Funcionario{}, err
	}
	if len(list) > 0 {
		res = list[0]
	} else {
		return res, model.ErrNotFound
	}
	return
}

func (p *pgUserRepository) fetchUserDetail(ctx context.Context, query string, args ...interface{}) (result []user.UserDetail, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]user.UserDetail, 0)
	for rows.Next() {
		t := user.UserDetail{}
		err = rows.Scan(
			&t.Id,
			&t.Nombre,
			&t.Apellido,
			&t.Celular,
			&t.Email,
			&t.SuperiorId,
			&t.EmpresaId,
			&t.Telefono,
			&t.CreatedOn,
			&t.UpdatedOn,
			&t.UserId,
			&t.Estado,
			&t.ProfilePhoto,
			&t.Rol,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}

func (p *pgUserRepository) fetchFuncionarios(ctx context.Context, query string, args ...interface{}) (result []user.Funcionario, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]user.Funcionario, 0)
	for rows.Next() {
		t := user.Funcionario{}
		err = rows.Scan(
			&t.FuncionarioId,
			&t.Nombre,
			&t.Apellido,
			&t.Celular,
			&t.Email,
			&t.SuperiorId,
			&t.EmpresaId,
			&t.Telefono,
			&t.CreatedOn,
			&t.UpdatedOn,
			&t.UserId,
			&t.Areas,
			&t.Estado,
			&t.ProfilePhoto,
			&t.Rol,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}

func (p *pgUserRepository) fetchClientes(ctx context.Context, query string, args ...interface{}) (result []user.Cliente, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]user.Cliente, 0)
	for rows.Next() {
		t := user.Cliente{}
		err = rows.Scan(
			&t.Id,
			&t.Nombre,
			&t.Apellido,
			&t.Celular,
			&t.Email,
			&t.SuperiorId,
			&t.EmpresaId,
			&t.Telefono,
			&t.CreatedOn,
			&t.UpdatedOn,
			&t.UserId,
			&t.Areas,
			&t.Estado,
			&t.ProfilePhoto,
			&t.IsAdmin,
			&t.Rol,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}

func (p *pgUserRepository) fetchUserShortInfo(ctx context.Context, query string, args ...interface{}) (result []user.UserShortInfo, err error) {
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]user.UserShortInfo, 0)
	for rows.Next() {
		t := user.UserShortInfo{}
		err = rows.Scan(
			&t.Id,
			&t.Nombre,
			&t.Apellido,
			&t.IsAdmin,
			&t.Pendiente,
			&t.Email,
			&t.Photo,
			&t.Estado,
			&t.DateTime,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}


func (p *pgUserRepository) fetchUsersForList(ctx context.Context, query string, args ...interface{}) (result []user.UserForList, err error) {
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

func (p *pgUserRepository) fetchUserArea(ctx context.Context, query string, args ...interface{}) (result []user.UserForList, err error) {
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
