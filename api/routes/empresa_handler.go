package routes

import (
	"log"
	"net/http"
	"soporte-go/core/model/empresa"

	// "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type EmpresaHandler struct {
	EmpresaUseCase empresa.EmpresaUseCase
}

func NewEmoresaHandler(e *echo.Echo, u empresa.EmpresaUseCase) {
	handler := &EmpresaHandler{
		EmpresaUseCase: u,
	}

	e.GET("/empresa/", handler.GetEmpresaUser)
	e.GET("/empresa/areas/", handler.GetAreasUserAdmin)
	e.GET("/empresa/areas-user/", handler.GetAreasUser)

	e.GET("/empresa/areas/:areaName/", handler.GetAreaDetail)
	e.POST("/empresa/", handler.StoreEmpresa)
	e.POST("/empresa/create-area/", handler.StoreArea)
	e.POST("/empresa/add-user-to-area/",handler.AddUserToArea)
}			
func (u *EmpresaHandler) AddUserToArea(c echo.Context) (err error){
	token := c.Request().Header["Authorization"][0]
	_,err = ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	var addUserRequestData empresa.AddUserRequestData;
	err = c.Bind(&addUserRequestData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	err = u.EmpresaUseCase.AddUserToArea(ctx,addUserRequestData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	log.Println(err)
	return c.JSON(http.StatusOK,nil)
	// userIds := c.FormValue("ids")
	// ctx := c.Request().Context()
	// err := u.EmpresaUseCase.AddUser/ToArea(ctx)
}

func (u *EmpresaHandler) GetAreaDetail(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_,err = ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	n := c.Param("areaName")
	ctx := c.Request().Context()
	res,err := u.EmpresaUseCase.GetAreaByName(ctx,n)
	// empresa.Id = casoId
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) StoreArea(c echo.Context) (err error) {
	var area empresa.Area
	token := c.Request().Header["Authorization"][0]
	claims,err := ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	err = c.Bind(&area)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}
	area.CreadorId = claims.UserId
	area.EmpresaId = claims.Empresa
	ctx := c.Request().Context()
	id,err := u.EmpresaUseCase.StoreArea(ctx,&area)
	area.Id = id
	// empresa.Id = casoId
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, area)
}

func (u *EmpresaHandler) GetAreasUserAdmin(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetAreasUserAdmin(ctx, claims.UserId)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) GetAreasUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetAreasUser(ctx, claims.UserId)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) GetEmpresaUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetEmpresa(ctx, claims.UserId, claims.Rol)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}



func (u *EmpresaHandler) StoreEmpresa(c echo.Context) (err error) {
	var empresa empresa.Empresa
	// token := c.Request().Header["Authorization"][0]
	// userId,err := ExtractClaims(token)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	// }
	err = c.Bind(&empresa)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	err = u.EmpresaUseCase.StoreEmpresa(ctx, &empresa)
	log.Println(err)
	// empresa.Id = casoId
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, empresa)

}
