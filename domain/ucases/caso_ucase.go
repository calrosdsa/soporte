package ucases

import (
	"context"
	"log"
	"soporte-go/core/model"
	"soporte-go/core/model/caso"

	"time"
)

type casoUseCase struct {
	casoRepo       caso.CasoRepository
	contextTimeout time.Duration
	util           model.Util
}

func NewCasoUseCase(uc caso.CasoRepository, timeout time.Duration, util model.Util) caso.CasoUseCase {
	return &casoUseCase{
		casoRepo:       uc,
		contextTimeout: timeout,
		util:           util,
	}
}

func (uc *casoUseCase) AsignarFuncionario(ctx context.Context, id string, idF string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.casoRepo.AsignarFuncionario(ctx, id, idF)
	return
}

func (uc *casoUseCase) GetCaso(ctx context.Context, id string) (res caso.Caso, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.casoRepo.GetCaso(ctx, id)
	if err != nil {
		return
	}
	return
}

func (uc *casoUseCase) GetCasosUser(ctx context.Context, id string, query *caso.CasoQuery, rol int) (res []caso.Caso, size int, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	page, offset := uc.util.PaginationValues(query.Page, query.PageSize)
	query.Page = page
	query.PageSize = offset
	// log.Println(query.Page,query.PageSize)
	if uc.util.IsClienteRol(rol) {
		res, err = uc.casoRepo.GetCasosCliente(ctx, id, query)
		if err != nil {
			return
		}
		size, err = uc.casoRepo.GetCasosCountCliente(ctx, id)
	} else if uc.util.IsFuncionarioRol(rol) {
		res, err = uc.casoRepo.GetCasosFuncionario(ctx, id, query)
		if err != nil {
			return
		}
		size, err = uc.casoRepo.GetCasosCountFuncionario(ctx, id)
	}
	return
}

func (uc *casoUseCase) GetAllCasosUser(ctx context.Context, id string, query *caso.CasoQuery) (res []caso.Caso, size int, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	size, err = uc.casoRepo.GetCasosCountbySuperiorId(ctx, id)
	if err != nil {
		return
	}
	res, err = uc.casoRepo.GetAllCasosUserCliente(ctx, id, query)
	return
}

func (uc *casoUseCase) UpdateCaso(ctx context.Context, columns []string, values ...interface{}) error {
	return nil
}

func (uc *casoUseCase) StoreCaso(ctx context.Context, cas *caso.Caso, id string, emI int) (idCaso string, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	idCaso, err = uc.casoRepo.StoreCaso(ctx, cas, id, emI)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
