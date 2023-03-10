package repository

import (
	"context"
	"errors"
	"reflect"
	"time"

	// "database/sqP"executing query""
	"log"
	"soporte-go/core/model"
	account "soporte-go/core/model/account"
	"soporte-go/core/model/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)
type Role byte
const (
	RoleCliente Role = iota
	RoleClienteAdmin
	RoleFuncionario
	RoleFuncionarioAdmin
)

type Estado byte
const (
	Activo Estado = iota
	Inactivo 
	Eliminado
)

type pgAccountRepository struct {
	Conn    *pgxpool.Pool
	Context context.Context
}

func NewPgAccountRepository(conn *pgxpool.Pool, ctx context.Context) account.AccountRepository {
	return &pgAccountRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func (p *pgAccountRepository) DeleteUser(ctx context.Context,id string)(err error){
	// query := `update users set estado = $1 where user_id = $2;`
	// query2 := `update clientes set estado = $1 where user_id = $2;`
	// _,err = p.Conn.Exec(ctx,query,Eliminado,id)
	// _,err = p.Conn.Exec(ctx,query2,Eliminado,id)
	query := `delete from users where user_id = $1;`
	query2 := `delete from clientes where user_id = $1;`
	_,err = p.Conn.Exec(ctx,query,id)
	_,err = p.Conn.Exec(ctx,query2,id)
	return
}

func (p *pgAccountRepository) Login(ctx context.Context, loginRequest *account.LoginRequest) (res user.ClienteAuth, err error) {
	var userId string
	query := `select user_id from users where email = $1 and password = crypt($2, password);`
	err = p.Conn.QueryRow(ctx,query,loginRequest.Email,loginRequest.Password).Scan(&userId)
	if err != nil{
		return res,model.ErrNotFound
	}
	res = user.ClienteAuth{}
	query1 := `select (client_id,email,estado,rol,empresa_id,nombre) from clientes where user_id = $1 limit 1;`
	err = p.Conn.QueryRow(ctx,query1,userId).Scan(&res)
	if err != nil{
		query1 := `select (funcionario_id,email,estado,rol,empresa_id,nombre) from funcionarios where user_id = $1 limit 1;`
	    err = p.Conn.QueryRow(ctx,query1,userId).Scan(&res)
	    if err != nil{
		    return res,model.ErrNotFound
	    }
		// return res,model.ErrNotFound
	}
	return
	// list, err := p.fetch(ctx, query, loginRequest.Email, loginRequest.Password)
	// if err != nil {
	// 	return account.User{}, err
	// }
	// if len(list) > 0 {
	// 	res = list[0]
	// } else {
	// 	return res, model.ErrNotFound
	// }
}

func (p *pgAccountRepository) RegisterCliente(ctx context.Context, a *account.RegisterForm) (res user.ClienteResponse, err error) {
	conn, err := p.Conn.BeginTx(p.Context, pgx.TxOptions{})

	if err != nil {
		return user.ClienteResponse{}, err
	}
	defer func() {
		if err != nil {
			conn.Rollback(context.TODO())
		} else {
			conn.Commit(context.TODO())
		}
	}()
	var emailInvitacion string
	query0 :=`select email from invitaciones where email = $1;`
	err = conn.QueryRow(ctx,query0,a.Email).Scan(&emailInvitacion)
	if emailInvitacion != *a.Email{
		return res,errors.New("No cuentas con una invitacion para poder registrarte")
	}

	t := account.User{}
	query := `select email,username from users where username = $1 or email = $2;`
	err = conn.QueryRow(p.Context, query, a.Username, a.Email).Scan(&t.Email, &t.Username)
	log.Println(t)
	if t.Email != nil {
		if *t.Email == *a.Email {
			return user.ClienteResponse{}, errors.New("Ya existe un usuario con este email")
		}
	}
	if t.Username != nil {
		if *t.Username == *a.Username {
			return user.ClienteResponse{}, errors.New("Este nombre ya esta ocupado")
		}
	}
	query2 := `insert into users (email,username,created_on,password) values ($1,$2,now(),crypt($3, gen_salt('bf'))) returning (user_id);`
	var userId string
	err = conn.QueryRow(p.Context, query2, a.Email, a.Username, a.Password).Scan(&userId)
	if err != nil {
		return user.ClienteResponse{}, err
	}
	cliente := user.ClienteResponse{}
	log.Println(reflect.TypeOf(a.EmpresaId))
	query3 := `insert into clientes (nombre,email,empresa_id,created_on,user_id,superior_id,rol) values ($1,$2,$3,$4,$5,$6,$7)
	returning (client_id,nombre,email,empresa_id,user_id);`
	err = conn.QueryRow(p.Context, query3, a.Username, a.Email, a.EmpresaId, time.Now(), userId,a.SuperiorId,a.Rol).Scan(&cliente)
	if err != nil {
		return user.ClienteResponse{}, err
	}
	query4 := `delete from invitaciones where email = $1;`
	_,err = conn.Exec(ctx,query4,cliente.Email)
	if err != nil{
		return user.ClienteResponse{},err
	}
	return cliente, err
}

func (p *pgAccountRepository) RegisterFuncionario(ctx context.Context, a *account.RegisterForm,id string) (user.FuncionarioResponse, error) {
	conn, err := p.Conn.BeginTx(p.Context, pgx.TxOptions{})
	if err != nil {
		return user.FuncionarioResponse{}, err
	}
	defer func() {
		if err != nil {
			conn.Rollback(context.TODO())
		} else {
			conn.Commit(context.TODO())
		}
	}()
	t := account.User{}
	query := `select email,username from users where username = $1 or email = $2;`
	err = conn.QueryRow(p.Context, query, a.Username, a.Email).Scan(&t.Email, &t.Username)
	log.Println(t)
	if t.Email != nil {
		if *t.Email == *a.Email {
			return user.FuncionarioResponse{}, errors.New("Ya existe un usuario con este email")
		}
	}
	if t.Username != nil {
		if *t.Username == *a.Username {
			return user.FuncionarioResponse{}, errors.New("Este nombre ya esta ocupado")
		}
	}
	query2 := `insert into users (email,username,created_on,password,rol) values ($1,$2,$3,crypt($4, gen_salt('bf')),$5) returning (user_id);`
	var userId string
	err = conn.QueryRow(p.Context, query2, a.Email, a.Username, time.Now(), a.Password, a.Rol).Scan(&userId)
	if err != nil {
		return user.FuncionarioResponse{}, err
	}
	funcionario := user.FuncionarioResponse{}
	query3 := `insert into funcionarios (nombre,email,empresa_id,created_on,user_id,superior_id) values ($1,$2,$3,$4,$5,$6)
	returning (funcionario_id,nombre,email,empresa_id,estado,created_on,user_id);`
	err = conn.QueryRow(p.Context, query3, a.Username, a.Email, a.EmpresaId, time.Now(), userId,id).Scan(&funcionario)
	if err != nil {
		return user.FuncionarioResponse{}, err
	}
	return funcionario, err
}

func (p *pgAccountRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []account.User, err error) {
	rows, err := p.Conn.Query(p.Context, query, args...)
	defer func ()  {
		rows.Close()
	   }()
	result = make([]account.User, 0)
	t := account.User{}
	for rows.Next() {
		err = rows.Scan(
			&t.UserId,
			&t.Username,
			&t.Password,
			&t.LastLogin,
			&t.CreatedOn,
			&t.Email,
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
