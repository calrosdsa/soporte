package ucases

import (
	"bytes"
	"context"
	"log"
	"soporte-go/core/model"
	"soporte-go/core/model/caso"
	"soporte-go/core/reportes/excel"
	"soporte-go/core/reportes/pdf"
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

func (uc *casoUseCase) GetReporteCasos(ctx context.Context, t model.FileType,options *caso.CasoReporteOptions) (b bytes.Buffer, err error) {
	var buffer bytes.Buffer
	casos, err := uc.casoRepo.GetCasosCliForReporte(ctx,options)
	casos2, err := uc.casoRepo.GetCasosFunForReporte(ctx,options)
	if err != nil {
		log.Println(err)
		return 
	}
	log.Println(casos)
	switch t {
		case model.XLSX:
			err = excel.ReporteCasosExcel(casos,casos2, &buffer)
			if err != nil {
				log.Println(err)
				return
			}
	    case model.PDF:
			err = pdf.ReporteCasos(casos, &buffer)
			if err != nil {
				log.Println(err)
				return
			}
	}
	return buffer, err
}

func (uc *casoUseCase) FinalizarCaso(ctx context.Context, fD *caso.FinalizacionDetail) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.casoRepo.FinalizarCaso(ctx, fD)
	return
}

func (uc *casoUseCase) AsignarFuncionario(ctx context.Context, id string, idF string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.casoRepo.AsignarFuncionario(ctx, id, idF)
	return
}

func (uc *casoUseCase) GetCaso(ctx context.Context, id string,rol int) (res caso.Caso, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	if uc.util.IsClienteRol(rol) {
		res, err = uc.casoRepo.GetCasoCliente(ctx, id)
		if err != nil {
			return
		}
	}else if uc.util.IsFuncionarioRol(rol){
		res, err = uc.casoRepo.GetCasoFuncionario(ctx, id)
		if err != nil {
			return
		}
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

func (uc *casoUseCase) GetAllCasosUser(ctx context.Context, id string, query *caso.CasoQuery, rol int) (res []caso.Caso, size int, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	page, offset := uc.util.PaginationValues(query.Page, 30)
	query.Page = page
	query.PageSize = offset
	if uc.util.IsClienteAdmin(rol) {
		size, err = uc.casoRepo.GetCasosCountbySuperiorId(ctx, id)
		if err != nil {
			return
		}
		res, err = uc.casoRepo.GetAllCasosUserCliente(ctx, id, query)
	} else if uc.util.IsFuncionarioAdmin(rol) {
		size, err = uc.casoRepo.GetCasosCount(ctx)
		if err != nil {
			return
		}
		res, err = uc.casoRepo.GetAllCasosUserFuncionario(ctx, 0, query)
	}
	return
}

func (uc *casoUseCase) UpdateCaso(ctx context.Context,c *caso.Caso) (err error) {
	ctx,cancel := context.WithTimeout(ctx,uc.contextTimeout)
	defer cancel()
	err = uc.casoRepo.UpdateCaso(ctx,c)
	return 
}

func (uc *casoUseCase) CreateCaso(ctx context.Context, cas *caso.Caso, id string, emI int,rol int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	log.Println(rol)
	if uc.util.IsClienteRol(rol){
		err = uc.casoRepo.CreateCasoCliente(ctx, cas, id, emI,rol)
		if err != nil {
			log.Println(err)
			return
		}
	}else if uc.util.IsFuncionarioRol(rol){
		err = uc.casoRepo.CreateCasoFuncionario(ctx,cas,id,emI,rol)
		if err != nil {
			return
		}
	}
	return
}
