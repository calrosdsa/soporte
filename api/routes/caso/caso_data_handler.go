package caso

import (
	"log"
	"net/http"
	"soporte-go/core/model"
	"soporte-go/core/model/caso"
	// "strconv"

	"github.com/labstack/echo/v4"
	_routes "soporte-go/api/routes"

)

type CasoDataHandler struct {
	CasoDataUseCase caso.CasoDataUseCase
}



func NewCasoDataHandler(e *echo.Echo,uc caso.CasoDataUseCase){
	handler := &CasoDataHandler{
		CasoDataUseCase: uc,
	}
	e.POST("/caso-data/casos-count/",handler.GetCountCasosByEstado)
	// e.GET("/caso-data/casos-count/:PId/",handler.GetCountCasosByEstado)
	e.POST("/caso-data/created-last-month/",handler.GetCasosCountCreatedLast30Days)
	e.GET("/caso-data/casos-count-proyectos/",handler.GetCasosData)
	e.POST("/caso-data/casos-estado-date/",handler.GetCasosEstadoByDate)
}

func (u *CasoDataHandler)GetCasosEstadoByDate (c echo.Context) (err error) {
	var data caso.ChartFilterData
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res,err := u.CasoDataUseCase.GetCasosEstadoByDate(ctx,data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK,res)
}

func (u *CasoDataHandler)GetCasosData(c echo.Context)(err error){
	// pId,err := strconv.Atoi(c.Param("PId"))
	token := c.Request().Header["Authorization"][0]
	claims, err := _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res,res2,res3,err := u.CasoDataUseCase.GetCasosData(ctx,claims.UserId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}
	response := caso.DataCasoResponse{
		ProyectosCasos:res,
		CasosLastMonth: res2,
		ProyectosCasosEstado: res3,
	}
	return c.JSON(http.StatusOK,response)
}

func (u *CasoDataHandler)GetCountCasosByEstado(c echo.Context)(err error){
	var data caso.ChartFilterData
	
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res,err := u.CasoDataUseCase.GetCasosCountEstadoByProyecto(ctx,data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK,res)
}

func (u *CasoDataHandler) GetCasosCountCreatedLast30Days(c echo.Context)(err error){
	// pId,err := strconv.Atoi(c.Param("PId"))
	var data caso.ChartFilterData
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res,err := u.CasoDataUseCase.GetCasosCountCreatedLast30Days(ctx,data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK,res)
}