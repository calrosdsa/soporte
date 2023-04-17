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
	// "github.com/sirupsen/logrus"
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

func (p *pgAccountRepository) DeleteUser(ctx context.Context, id string) (err error) {
	// query := `update users set estado = $1 where user_id = $2;`
	// query2 := `update clientes set estado = $1 where user_id = $2;`
	// _,err = p.Conn.Exec(ctx,query,Eliminado,id)
	// _,err = p.Conn.Exec(ctx,query2,Eliminado,id)
	query := `delete from users where user_id = $1;`
	query2 := `delete from clientes where user_id = $1;`
	_, err = p.Conn.Exec(ctx, query, id)
	if err != nil {
		return
	}
	_, err = p.Conn.Exec(ctx, query2, id)
	if err != nil {
		return
	}
	return
}

func (p *pgAccountRepository) Login(ctx context.Context, loginRequest *account.LoginRequest) (res user.UserAuth, err error) {
	var userId string
	// log.Println(loginRequest.Email)
	// log.Println(loginRequest.Password)
	query := `select user_id from users where email = $1 and password = crypt($2, password);`
	err = p.Conn.QueryRow(ctx, query, loginRequest.Email, loginRequest.Password).Scan(&userId)
	if err != nil {
		log.Println(err)
		// log.Println("error is here")
		return res, model.ErrNotFound
	}
	res = user.UserAuth{}
	query1 := `select (client_id,email,estado,rol,empresa_id,nombre) from clientes where user_id = $1 limit 1;`
	err = p.Conn.QueryRow(ctx, query1, userId).Scan(&res)
	if err != nil {
		log.Println("isFuncionario")
		query1 := `select (funcionario_id,email,estado,rol,empresa_id,nombre) from funcionarios where user_id = $1 limit 1;`
		err = p.Conn.QueryRow(ctx, query1, userId).Scan(&res)
		if err != nil {
			return res, model.ErrNotFound
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

func (p *pgAccountRepository) RegisterCliente(ctx context.Context, a *account.RegisterForm) (res user.UserAuth, err error) {
	conn, err := p.Conn.BeginTx(p.Context, pgx.TxOptions{})

	if err != nil {
		return user.UserAuth{}, err
	}
	defer func() {
		if err != nil {
			conn.Rollback(context.TODO())
		} else {
			conn.Commit(context.TODO())
		}
	}()
	var emailInvitacion string
	query0 := `select email from invitaciones where email = $1;`
	err = conn.QueryRow(ctx, query0, a.Email).Scan(&emailInvitacion)
	if emailInvitacion != a.Email {
		return res, errors.New("no cuentas con una invitacion para poder registrarte")
	}

	t := account.User{}
	query := `select email,username from users where username = $1 or email = $2;`
	err = conn.QueryRow(p.Context, query, a.Nombre, a.Email).Scan(&t.Email, &t.Username)
	log.Println(t)
	if t.Email != nil {
		if *t.Email == a.Email {
			return user.UserAuth{}, errors.New("ya existe un usuario con este email")
		}
	}
	if t.Username != nil {
		if *t.Username == a.Nombre {
			return user.UserAuth{}, errors.New("yste nombre ya esta ocupado")
		}
	}
	query2 := `insert into users (email,username,created_on,password) values ($1,$2,now(),crypt($3, gen_salt('bf')))
	returning (user_id);`
	var userId string
	err = conn.QueryRow(p.Context, query2, a.Email, a.Nombre, a.Password).Scan(&userId)
	if err != nil {
		return user.UserAuth{}, err
	}
	cliente := user.UserAuth{}
	log.Println(reflect.TypeOf(a.EmpresaId))
	query3 := `insert into clientes (nombre,apellido,email,empresa_id,created_on,user_id,rol,superior_id) values ($1,$2,$3,$4,$5,$6,$7,$8)
	returning (client_id,email,estado,rol,empresa_id,(''));`
	err = conn.QueryRow(p.Context, query3, a.Nombre,a.Apellido ,a.Email, a.EmpresaId, time.Now(), userId, a.Rol, a.SuperiorId).Scan(&cliente)
	if err != nil {
		return user.UserAuth{}, err
	}
	query4 := `delete from invitaciones where email = $1;`
	_, err = conn.Exec(ctx, query4, cliente.Email)
	if err != nil {
		return user.UserAuth{}, err
	}
	return cliente, err
}

func (p *pgAccountRepository) RegisterFuncionario(ctx context.Context, a *account.RegisterForm) (res user.UserAuth, err error) {
	conn, err := p.Conn.BeginTx(p.Context, pgx.TxOptions{})
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			conn.Rollback(context.TODO())
		} else {
			conn.Commit(context.TODO())
		}
	}()
	var query string
	t := account.User{}
	query = `select email,username from users where username = $1 or email = $2;`
	err = conn.QueryRow(p.Context, query, a.Nombre, a.Email).Scan(&t.Email, &t.Username)
	// log.Println(t)
	if t.Email != nil {
		if *t.Email == a.Email {
			return user.UserAuth{}, model.ErrConflictEmail
		}
	}
	if t.Username != nil {
		if *t.Username == a.Nombre {
			return user.UserAuth{}, model.ErrConflictUsername
		}
	}
	query = `insert into users (email,username,created_on,password) values ($1,$2,$3,crypt($4, gen_salt('bf'))) returning (user_id);`
	var userId string
	err = conn.QueryRow(p.Context, query, a.Email, a.Nombre, time.Now(), a.Password).Scan(&userId)
	if err != nil {
		return
	}
	// var superiorId string
	// log.Println(a.SuperiorId)
	// query = `select superior_id from funcionarios where funcionario_id = $1;`
	// err = conn.QueryRow(p.Context, query, a.SuperiorId).Scan(&superiorId)
	// if err != nil {
	// 	log.Println("error is here")
	// 	return
	// }
	// log.Panicln("USER INSERTED")
	res = user.UserAuth{}
	query = `insert into funcionarios (nombre,apellido,email,empresa_id,created_on,user_id,rol,superior_id) values ($1,$2,$3,$4,$5,$6,$7,$8)
	returning (funcionario_id,email,estado,rol,empresa_id,(''));`
	err = conn.QueryRow(p.Context, query, a.Nombre,a.Apellido, a.Email, a.EmpresaId, time.Now(), userId, a.Rol, a.SuperiorId).Scan(&res)
	// log.Println(*t.Username)
	if err != nil {
		return
	}
	query = `delete from invitaciones where email = $1;`
	_, err = conn.Exec(p.Context, query, res.Email)
	if err != nil {
		return
	}

	return
}

// func (p *pgAccountRepository) ValidateInvitation(ctx context.Context,mail *string,rol *int)(err error){
// 	var query string
// 	var email string
// 	if *rol == int(model.RoleCliente) || *rol == int(model.RoleClienteAdmin) {
// 		query = `select email from invitaciones where email = $1;`
// 		err = p.Conn.QueryRow(ctx,query,mail).Scan(&email)
// 	}else if *rol == int(RoleFuncionario) || *rol == int(RoleFuncionarioAdmin) {

// 	}
// }

// func (p *pgAccountRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []account.User, err error) {
// 	rows, err := p.Conn.Query(p.Context, query, args...)
// 	defer func ()  {
// 		rows.Close()
// 	   }()
// 	result = make([]account.User, 0)
// 	t := account.User{}
// 	for rows.Next() {
// 		err = rows.Scan(
// 			&t.UserId,
// 			&t.Username,
// 			&t.Password,
// 			&t.LastLogin,
// 			&t.CreatedOn,
// 			&t.Email,
// 			&t.Estado,
// 		)
// 		if err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		result = append(result, t)
// 		log.Println(result)
// 	}

// 	return result, nil
// }
