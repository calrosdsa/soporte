package user
import (
	"context"
)


type UserRegistrationRequest struct {
	To []string `json:"to"`
	IsAdmin bool `json:"is_admin"`
}


type UserUseCases interface{
	// CreateCliente(ctx context.Context,user *Cliente ) (res string,err error)
	UserRegisterInvitation(ctx context.Context, url string, to UserRegistrationRequest,id string) ([]UserShortInfo,error)
	GetClientesByArea(context.Context,int)([]UserArea,error)
	UpdateCliente(ctx context.Context,columns []string,values ...interface{}) (error)
	UpdateFuncionario(ctx context.Context,columns []string,values ...interface{}) (error)
	GetClienteById(ctx context.Context,id string) (Cliente,error)
	GetClientes(ctx context.Context,id string) ([]UserShortInfo,error)
	GetFuncionarios(ctx context.Context) ([]Funcionario,error)
	GetFuncionarioById(ctx context.Context,id string) (Funcionario,error)
	GetUsersShortIInfo(ctx context.Context,id string) ([]UserShortInfo,error)
	ValidateEmail(ctx context.Context,email string)(error)
	ReSendEmail(m []string,url string)
	DeleteInvitation(ctx context.Context ,m string)(err error)
	GetClientesFiltered(ctx context.Context,f int) ([]UserArea,error)
	SearchUser(ctx context.Context,id string,q string)([]UserShortInfo,error)

}


type UserRepository interface{
	GetUsersShortIInfo(ctx context.Context,id string) ([]UserShortInfo,error)
	GetInvitaciones(ctx context.Context,id string) ([]UserShortInfo,error)
	CreateUserInvitation(context.Context,*UserShortInfo) (UserShortInfo,error)
	// CreateCliente(ctx context.Context,user *Cliente ) (res string,err error)
	GetClientesFiltered(ctx context.Context,f int) ([]UserArea,error)
	GetClientesByArea(context.Context,int)([]UserArea,error)
	UpdateCliente(ctx context.Context,columns []string,values ...interface{}) (error)
	UpdateFuncionario(ctx context.Context,columns []string,values ...interface{}) (error)
	GetClienteById(ctx context.Context,id string) (res Cliente,err error)
	GetClientes(ctx context.Context) (clientes []Cliente,err error)
	GetFuncionarios(ctx context.Context) (funcionarios []Funcionario,err error)
	GetFuncionarioById(ctx context.Context,id string) (res Funcionario,err error)
	ValidateEmail(ctx context.Context,m string)(error)
	DeleteInvitation(ctx context.Context ,m string)(err error)
	SearchUser(ctx context.Context,id string,q string)([]UserShortInfo,error)


}