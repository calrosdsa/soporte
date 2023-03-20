package caso

import (
	"log"
	"net/http"
	"reflect"
	"soporte-go/core/model/caso"
	model "soporte-go/core/model"
	"strconv"

	// "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	_routes "soporte-go/api/routes"
	"github.com/sirupsen/logrus"

)

type CasoHandler struct {
	CasoUseCase caso.CasoUseCase
}

func NewCasoHandler(e *echo.Echo, uc caso.CasoUseCase) {
	handler := &CasoHandler{
		CasoUseCase: uc,
	}

	e.POST("/caso/", handler.StoreCaso)
	e.GET("/caso/:casoId/", handler.GetCaso)
	e.GET("/casos", handler.GetCasosUser)
	e.GET("/casos-all", handler.GetAllCasosUser)
	
	e.GET("/caso/asignar-funcionario/:idCaso/:idFuncionario/",handler.AsignarFuncionario)
}

func (u *CasoHandler) AsignarFuncionario(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	id := c.Param("idCaso")
	idF := c.Param("idFuncionario")
	ctx := c.Request().Context()
	err = u.CasoUseCase.AsignarFuncionario(ctx,id,idF)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK,nil)
}

func (u *CasoHandler) GetAllCasosUser(c echo.Context) (err error) {
	// page := c.QueryParam("page")
	estado := c.QueryParam("estado")
	prioridad := c.QueryParam("prioridad")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	log.Println("pagevalue", page)
	casoQuery := caso.CasoQuery{
		Page:      page,
		Estado:    estado,
		Prioridad: prioridad,
	}
	priori, err := strconv.Atoi(prioridad)
	log.Println(err)
	log.Println(reflect.TypeOf(priori))
	log.Println(casoQuery)
	token := c.Request().Header["Authorization"][0]
	claims, err := _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, size, err := u.CasoUseCase.GetAllCasosUser(ctx, claims.UserId, &casoQuery)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	response := caso.CasosResponse{
		Casos:   res,
		Size:    size,
		Current: page,
	}
	return c.JSON(http.StatusOK, response)
}

func (u *CasoHandler) GetCasosUser(c echo.Context) (err error) {
	// page := c.QueryParam("page")
	estado := c.QueryParam("estado")
	prioridad := c.QueryParam("prioridad")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	log.Println("pagevalue", page)
	casoQuery := caso.CasoQuery{
		Page:      page,
		Estado:    estado,
		Prioridad: prioridad,
	}
	// priori, err := strconv.Atoi(prioridad)
	// log.Println(err)
	// log.Println(reflect.TypeOf(priori))
	// log.Println(casoQuery)
	token := c.Request().Header["Authorization"][0]
	claims, err := _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, size, err := u.CasoUseCase.GetCasosUser(ctx, &claims.UserId, &casoQuery,&claims.Rol)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	response := caso.CasosResponse{
		Casos:   res,
		Size:    size,
		Current: page,
	}
	return c.JSON(http.StatusOK, response)
}

func (u *CasoHandler) GetCaso(c echo.Context) (err error) {
	id := c.Param("casoId")

	// token := c.Request().Header["Authorization"][0]
	// userId,err := ExtractClaims(token)
	ctx := c.Request().Context()
	res, err := u.CasoUseCase.GetCaso(ctx, id)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *CasoHandler) StoreCaso(c echo.Context) (err error) {
	var caso caso.Caso
	token := c.Request().Header["Authorization"][0]
	claims, err := _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,model.ResponseError{Message: err.Error()})
	}
	err = c.Bind(&caso)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	casoId, err := u.CasoUseCase.StoreCaso(ctx, &caso, claims.UserId, claims.Empresa)
	caso.Id = casoId
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, caso)
}


