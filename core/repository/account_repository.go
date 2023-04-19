package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"time"

	// "database/sqP"executing query""
	"log"
	"soporte-go/core/model"
	account "soporte-go/core/model/account"
	"soporte-go/core/model/user"

	// "github.com/sirupsen/logrus"
)

type pgAccountRepository struct {
	Conn    *sql.DB
	Context context.Context
}

func NewPgAccountRepository(conn *sql.DB, ctx context.Context) account.AccountRepository {
	return &pgAccountRepository{
		Conn:    conn,
		Context: ctx,
	}
}

func (p *pgAccountRepository) DeleteUser(ctx context.Context, id string) (err error) {
	// query := `update users set estado = $1 where user_id = $2;`
	// query2 := `update clientes set estado = $1 where user_id = $2;`
	// _,err = p.Conn.ExecContext(ctx,query,Eliminado,id)
	// _,err = p.Conn.ExecContext(ctx,query2,Eliminado,id)
	query := `delete from users where user_id = $1;`
	query2 := `delete from clientes where user_id = $1;`
	_, err = p.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return
	}
	_, err = p.Conn.ExecContext(ctx, query2, id)
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
	err = p.Conn.QueryRowContext(ctx, query, loginRequest.Email, loginRequest.Password).Scan(&userId)
	if err != nil {
		log.Println(err)
		// log.Println("error is here")
		return res, model.ErrNotFound
	}
	res = user.UserAuth{}
	query1 := `select client_id,email,estado,rol,empresa_id,nombre from clientes where user_id = $1 limit 1;`
	err = p.Conn.QueryRowContext(ctx, query1, userId).Scan(&res.Id,&res.Email,&res.Estado,&res.Rol,&res.EmpresaId,&res.Username)
	if err != nil {
		log.Println("isFuncionario")
		query1 := `select funcionario_id,email,estado,rol,empresa_id,nombre from funcionarios where user_id = $1 limit 1;`
		err = p.Conn.QueryRowContext(ctx, query1, userId).Scan(&res.Id,&res.Email,&res.Estado,&res.Rol,&res.EmpresaId,&res.Username)
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
	conn, err := p.Conn.BeginTx(p.Context, &sql.TxOptions{})

	if err != nil {
		return user.UserAuth{}, err
	}
	defer func() {
		if err != nil {
			conn.Rollback()
		} else {
			conn.Commit()
		}
	}()
	var emailInvitacion string
	query0 := `select email from invitaciones where email = $1;`
	err = conn.QueryRowContext(ctx, query0, a.Email).Scan(&emailInvitacion)
	if emailInvitacion != a.Email {
		return res, errors.New("no cuentas con una invitacion para poder registrarte")
	}

	t := account.User{}
	query := `select email,username from users where username = $1 or email = $2;`
	err = conn.QueryRowContext(p.Context, query, a.Nombre, a.Email).Scan(&t.Email, &t.Username)
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
	returning user_id;`
	var userId string
	err = conn.QueryRowContext(p.Context, query2, a.Email, a.Nombre, a.Password).Scan(&userId)
	if err != nil {
		return user.UserAuth{}, err
	}
	cliente := user.UserAuth{}
	log.Println(reflect.TypeOf(a.EmpresaId))
	query3 := `insert into clientes nombre,apellido,email,empresa_id,created_on,user_id,rol,superior_id values ($1,$2,$3,$4,$5,$6,$7,$8)
	returning client_id,email,estado,rol,empresa_id;`
	err = conn.QueryRowContext(p.Context, query3, a.Nombre,a.Apellido ,a.Email, a.EmpresaId, time.Now(), userId, a.Rol, a.SuperiorId).Scan(
		&cliente.Id,&cliente.Email,&cliente.EmpresaId,&cliente.EmpresaId,
	)
	if err != nil {
		return user.UserAuth{}, err
	}
	query4 := `delete from invitaciones where email = $1;`
	_, err = conn.ExecContext(ctx, query4, cliente.Email)
	if err != nil {
		return user.UserAuth{}, err
	}
	return cliente, err
}

func (p *pgAccountRepository) RegisterFuncionario(ctx context.Context, a *account.RegisterForm) (res user.UserAuth, err error) {
	conn, err := p.Conn.BeginTx(p.Context, &sql.TxOptions{})

	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			conn.Rollback()
		} else {
			conn.Commit()
		}
	}()
	var query string
	t := account.User{}
	query = `select email,username from users where username = $1 or email = $2;`
	err = conn.QueryRowContext(p.Context, query, a.Nombre, a.Email).Scan(&t.Email, &t.Username)
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
	err = conn.QueryRowContext(p.Context, query, a.Email, a.Nombre, time.Now(), a.Password).Scan(&userId)
	if err != nil {
		return
	}
	// var superiorId string
	// log.Println(a.SuperiorId)
	// query = `select superior_id from funcionarios where funcionario_id = $1;`
	// err = conn.QueryRowContext(p.Context, query, a.SuperiorId).Scan(&superiorId)
	// if err != nil {
	// 	log.Println("error is here")
	// 	return
	// }
	// log.Panicln("USER INSERTED")
	res = user.UserAuth{}
	query = `insert into funcionarios (nombre,apellido,email,empresa_id,created_on,user_id,rol,superior_id) values ($1,$2,$3,$4,$5,$6,$7,$8)
	returning funcionario_id,email,estado,rol,empresa_id;`
	err = conn.QueryRowContext(p.Context, query, a.Nombre,a.Apellido, a.Email, a.EmpresaId, time.Now(), userId, a.Rol, a.SuperiorId).Scan(
		&res.Id,&res.Email,&res.EmpresaId,&res.EmpresaId,
	)
	// log.Println(*t.Username)
	if err != nil {
		return
	}
	query = `delete from invitaciones where email = $1;`
	_, err = conn.ExecContext(p.Context, query, res.Email)
	if err != nil {
		return
	}

	return
}
