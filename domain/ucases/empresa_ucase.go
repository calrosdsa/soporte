package ucases

import (
	"context"
	// "log"
	"soporte-go/core/model"
	"soporte-go/core/model/empresa"
	"time"
)

type empresaUseCase struct {
	empresaRepo    empresa.EmpresaRepository
	contextTimeout time.Duration
	util           model.Util
}

func NewEmpresaUseCase(uc empresa.EmpresaRepository, timeout time.Duration, util model.Util) empresa.EmpresaUseCase {
	return &empresaUseCase{
		empresaRepo:    uc,
		contextTimeout: timeout,
		util:           util,
	}
}

func (uc *empresaUseCase) GetAreasEmpresa(ctx context.Context, emId int) (res []empresa.Area, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.empresaRepo.GetAreasEmpresa(ctx, emId)
	return
}

func (uc *empresaUseCase) GetSubAreas(ctx context.Context, parentId int) (res []empresa.Area, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.empresaRepo.GetSubAreas(ctx, parentId)
	return
}

func (uc *empresaUseCase) CreateSubArea(ctx context.Context, a *empresa.SubArea) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.empresaRepo.CreateSubArea(ctx, a)
	return
}

func (uc *empresaUseCase) AreaChangeState(ctx context.Context, state int, id int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.empresaRepo.AreaChangeState(ctx, state, id)
	return err
}

func (uc *empresaUseCase) GetAreaByName(ctx context.Context, n string) (res empresa.Area, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.empresaRepo.GetAreaByName(ctx, n)
	return
}

func (uc *empresaUseCase) AddUserToArea(ctx context.Context, a *empresa.AddUserRequestData) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	for _, value := range a.Users {
		err = uc.empresaRepo.AddUserToArea(ctx, value.Id, a)
	}
	return
}

func (uc *empresaUseCase) GetAreasUser(ctx context.Context, userId string) (res []empresa.AreaUser, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.empresaRepo.GetAreasUser(ctx, userId)
	return
}

func (uc *empresaUseCase) GetUsersAreaByAreaId(ctx context.Context, areaId int) (res []empresa.UserArea, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.empresaRepo.GetUsersAreaByAreaId(ctx, areaId)
	return
}

func (uc *empresaUseCase) GetAreasFromUser(ctx context.Context, userId string, emId int, rol int) (res []empresa.Area, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	if uc.util.IsClienteAdmin(rol) {
		res, err = uc.empresaRepo.GetAreasClienteAdmin(ctx, userId)
		return
	} else if uc.util.IsFuncionarioRol(rol) {
		res, err = uc.empresaRepo.GetSubAreasFuncionario(ctx, emId)
		return
	}
	return res, model.ErrConflict
}

func (uc *empresaUseCase) StoreEmpresa(ctx context.Context, empresa *empresa.Empresa) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.empresaRepo.StoreEmpresa(ctx, empresa)
	return
}

func (uc *empresaUseCase) StoreArea(ctx context.Context, area *empresa.Area) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.empresaRepo.StoreArea(ctx, area)
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

func (uc *empresaUseCase) GetEmpresas(ctx context.Context, emId *int) (res []empresa.Empresa, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.empresaRepo.GetEmpresas(ctx, emId)
	if err != nil {
		return res, err
	}
	return
}
