package ucases

import (
	"context"
	"log"
	"soporte-go/core/model/caso"
	"soporte-go/core/model/empresa"
	"time"
)

type casoDataUseCase struct {
	casoDataRepo   caso.CasoDataRepository
	empresaRepo    empresa.EmpresaRepository
	contextTimeout time.Duration
}

func NewCasoDataUseCase(uc caso.CasoDataRepository, uc2 empresa.EmpresaRepository, timeout time.Duration) caso.CasoDataUseCase {
	return &casoDataUseCase{
		casoDataRepo:   uc,
		contextTimeout: timeout,
		empresaRepo:    uc2,
	}
}

func (c *casoDataUseCase) GetCasosData(ctx context.Context, idUser string) (res []caso.DataCaso, res2 []caso.DataCaso,
	res3 []caso.DataCaso, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	proyectos, err := c.empresaRepo.GetAreasUser(ctx, idUser)
	ids := make([]int, len(proyectos))
	for i, val := range proyectos {
		ids[i] = val.Id
	}
	to := time.Now().Format("2006-01-02")
	from := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	d := caso.ChartFilterData{
		Ids:      ids,
		ToDate:   to,
		FromDate: from,
		TypeDate: "day",
	}
	log.Println(ids)
	res, err = c.casoDataRepo.GetCasosCountByProyecto(ctx, d)
	if err != nil {
		log.Println(err)
	}
	res2, err = c.casoDataRepo.GetCasosCountCreatedLast30Days(ctx, d)
	if err != nil {
		log.Println(err)
	}
	res3, err = c.casoDataRepo.GetCasosCountEstadoByProyecto(ctx, d)
	if err != nil {
		log.Println(err)
	}
	return
}

func (c *casoDataUseCase) GetCasosCountEstadoByProyecto(ctx context.Context, d caso.ChartFilterData) (res []caso.DataCaso, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	res, err = c.casoDataRepo.GetCasosCountEstadoByProyecto(ctx, d)
	return
}

func (c *casoDataUseCase) GetCasosCountCreatedLast30Days(ctx context.Context, d caso.ChartFilterData) (res []caso.DataCaso, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	res, err = c.casoDataRepo.GetCasosCountCreatedLast30Days(ctx, d)
	return
}

func (c *casoDataUseCase) GetCasosEstadoByDate(ctx context.Context, d caso.ChartFilterData) (res []caso.DataCasoEstado, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	res,err = c.casoDataRepo.GetCasosEstadoByDate(ctx,d)
	return
}
