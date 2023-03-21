package user

import (
	"context"
	// "soporte-go/core/model/empresa"
	// "time"
)

type UserRegistrationRequest struct {
	To      []string `json:"to"`
	IsAdmin bool     `json:"is_admin"`
	EmpresaId int `json:"empresa_id"`
}

type UserUseCases interface {
	// CreateCliente(ctx context.Context,user *Cliente ) (res string,err error)
	UserRegisterInvitation(ctx context.Context, to *UserRegistrationRequest, id string, rol int, empresaId int) ([]UserShortInfo, error)
	GetClientesByArea(context.Context, int) ([]UserArea, error)
	UpdateCliente(ctx context.Context, columns []string, values ...interface{}) error
	UpdateFuncionario(ctx context.Context, columns []string, values ...interface{}) error
	GetUserById(ctx context.Context, id string, rol int) (Cliente, error)
	GetClientes(ctx context.Context, id string, rol int) ([]UserShortInfo, error)
	GetFuncionarios(ctx context.Context) ([]Funcionario, error)
	GetFuncionarioById(ctx context.Context, id string) (Funcionario, error)

	GetUsersShortIInfo(ctx context.Context, id string,rol int,emId int) ([]UserShortInfo, error)

	ValidateEmail(ctx context.Context, email string) error
	ReSendEmail(m []string, url string)
	DeleteInvitation(ctx context.Context, m string) (err error)
	GetUserAddList(ctx context.Context, f int,rol int,sId string) ([]UserArea, error)
	SearchUser(ctx context.Context, id string, q string) ([]UserShortInfo, error)

	GetUsersbyEmpresaId(ctx context.Context,emId int)([]UserShortInfo,error	)
}

type UserRepository interface {

	GetUsersShortIInfoC(ctx context.Context, id string) ([]UserShortInfo, error)
	GetUsersShortIInfoF(ctx context.Context, emID int) ([]UserShortInfo, error)

	// GetUsersShortIInfo(ctx context.Context, id string, rol int) ([]UserShortInfo, error)
	GetInvitaciones(ctx context.Context, id string) ([]UserShortInfo, error)
	CreateUserInvitation(ctx context.Context, us *UserShortInfo, rol int) (UserShortInfo, error)

	GetClientesEmpresa(ctx context.Context,emId int)([]UserShortInfo,error)
	// CreateCliente(ctx context.Context,user *Cliente ) (res string,err error)
	GetUserAddList(ctx context.Context, f int, rol int,sId string) ([]UserArea, error)
	GetClientesByArea(context.Context, int) ([]UserArea, error)
	UpdateCliente(ctx context.Context, columns []string, values ...interface{}) error
	UpdateFuncionario(ctx context.Context, columns []string, values ...interface{}) error
	GetUserById(ctx context.Context, id string, rol int) (res Cliente, err error)
	GetClientes(ctx context.Context) (clientes []Cliente, err error)
	GetFuncionarios(ctx context.Context) (funcionarios []Funcionario, err error)
	GetFuncionarioById(ctx context.Context, id string) (res Funcionario, err error)
	ValidateEmail(ctx context.Context, m string) error
	DeleteInvitation(ctx context.Context, m string) (err error)
	SearchUser(ctx context.Context, id string, q string) ([]UserShortInfo, error)
}
