package ucases

import (
	"context"
	
	"soporte-go/core/model/account"
	"soporte-go/core/model/user"
	"time"
)

type accountUseCase struct {
	accountRepo account.AccountRepository
	contextTimeout   time.Duration
}


func NewAccountUseCase(a account.AccountRepository, timeout time.Duration) account.AccountUseCase {
	return &accountUseCase{
		accountRepo:   a,
		contextTimeout:   timeout,
	}
}

func (a *accountUseCase) DeleteUser(ctx context.Context,id string)(err error){
	ctx,cancel := context.WithTimeout(ctx,a.contextTimeout)
	defer cancel()
	err =  a.accountRepo.DeleteUser(ctx,id)
	return
}

func (a *accountUseCase) Login(ctx context.Context,loginRequest *account.LoginRequest) (res user.ClienteAuth,err error){
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
    res, err = a.accountRepo.Login(ctx,loginRequest)
	if err != nil {
		return
	}
	return 
 }
 
 func (a *accountUseCase) RegisterCliente(ctx context.Context,form *account.RegisterForm) (res user.ClienteResponse,err error){
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	res,err = a.accountRepo.RegisterCliente(ctx,form)
	return
 }
 
 func (a *accountUseCase) RegisterFuncionario(ctx context.Context,form *account.RegisterForm,id string) (res user.FuncionarioResponse,err error){
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	res,err = a.accountRepo.RegisterFuncionario(ctx,form,id)
	return
 }