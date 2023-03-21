package ucases

import (
	"context"

	"soporte-go/core/model/account"
	"soporte-go/core/model/user"
	"time"
	"soporte-go/core/model"
)

type accountUseCase struct {
	accountRepo    account.AccountRepository
	contextTimeout time.Duration
	util     model.Util
}

func NewAccountUseCase(a account.AccountRepository, timeout time.Duration,util model.Util) account.AccountUseCase {
	return &accountUseCase{
		accountRepo:    a,
		contextTimeout: timeout,
		util: util,
	}
}

func (a *accountUseCase) RegisterUser(ctx context.Context,form *account.RegisterForm)(res user.UserAuth,err error){
	ctx,cancel := context.WithTimeout(ctx,a.contextTimeout)
	defer cancel()
	if a.util.IsClienteRol(form.Rol){
		res,err =  a.accountRepo.RegisterCliente(ctx, form)
	} else if a.util.IsFuncionarioRol(form.Rol){
		res,err =  a.accountRepo.RegisterFuncionario(ctx, form)
	}
	return
}

func (a *accountUseCase) DeleteUser(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	err = a.accountRepo.DeleteUser(ctx, id)
	return
}

func (a *accountUseCase) Login(ctx context.Context, loginRequest *account.LoginRequest) (res user.UserAuth, err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	res, err = a.accountRepo.Login(ctx, loginRequest)
	if err != nil {
		return
	}
	return
}

// func (a *accountUseCase) RegisterCliente(ctx context.Context, form *account.RegisterForm) (res user.ClienteResponse, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
// 	defer cancel()

// 	res, err = a.accountRepo.RegisterCliente(ctx, form)
// 	return
// }

// func (a *accountUseCase) RegisterFuncionario(ctx context.Context, form *account.RegisterForm) (res user.UserAuth, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
// 	defer cancel()

// 	// res ,err = 

// 	res, err = a.accountRepo.RegisterFuncionario(ctx, form)
// 	return
// }
