package ucases

import (
	"context"
	// "log"
	"soporte-go/core/model/empresa"
	"time"
)

type empresaUseCase struct {
	empresaRepo    empresa.EmpresaRepository
	contextTimeout time.Duration
}

func NewEmpresaUseCase(uc empresa.EmpresaRepository, timeout time.Duration) empresa.EmpresaUseCase {
	return &empresaUseCase{
		empresaRepo:    uc,
		contextTimeout: timeout,
	}
}

func (uc *empresaUseCase) AreaChangeState(ctx context.Context,state int,id int)(err error){
	ctx,cancel := context.WithTimeout(ctx,uc.contextTimeout)
	defer cancel()
	err = uc.empresaRepo.AreaChangeState(ctx,state,id)
	return err
}

func (uc *empresaUseCase) GetAreaByName(ctx context.Context,n string) (res empresa.Area,err error){
	ctx, cancel := context.WithTimeout(ctx,uc.contextTimeout)
	defer cancel()
	res,err = uc.empresaRepo.GetAreaByName(ctx,n)
	return 
}

func (uc *empresaUseCase) AddUserToArea(ctx context.Context,a *empresa.AddUserRequestData)(err error){
	ctx,cancel := context.WithTimeout(ctx,uc.contextTimeout)
	defer cancel()
	for _,value := range a.Users{
		err = uc.empresaRepo.AddUserToArea(ctx,&value.Id,&value.Name,a)
	}
	return
}

func (uc *empresaUseCase) GetAreasUser(ctx context.Context, userId string) (res []empresa.AreaUser, err error) {
	ctx,cancel := context.WithTimeout(ctx,uc.contextTimeout)
	defer cancel()
	res ,err = uc.empresaRepo.GetAreasUser(ctx,userId)
	return
}

func (uc *empresaUseCase) GetAreasUserAdmin(ctx context.Context, userId *string,rol *int) (res []empresa.Area, err error) {
	ctx,cancel := context.WithTimeout(ctx,uc.contextTimeout)
	defer cancel()
	res ,err = uc.empresaRepo.GetAreasUserAdmin(ctx,userId,rol)
	return
}

func (uc *empresaUseCase) StoreEmpresa(ctx context.Context, empresa *empresa.Empresa) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.empresaRepo.StoreEmpresa(ctx, empresa)
	return
}

func (uc *empresaUseCase) StoreArea(ctx context.Context, area *empresa.Area) (err error) {
	ctx,cancel := context.WithTimeout(ctx,uc.contextTimeout)
	defer cancel()
	err = uc.empresaRepo.StoreArea(ctx,area)
	return 
}

func (uc *empresaUseCase) GetEmpresa(ctx context.Context, userId string, rol int) (res empresa.Empresa, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.empresaRepo.GetEmpresa(ctx, userId, rol)
	if err != nil {
		return res, err
	}
	return
}

func (uc *empresaUseCase) GetEmpresas(ctx context.Context,emId *int) (res []empresa.Empresa, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.empresaRepo.GetEmpresas(ctx,emId)
	if err != nil {
		return res, err
	}
	return
}
