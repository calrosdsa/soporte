package routes

import (
	"log"
	"net/http"
	"soporte-go/core/model/empresa"
	"strconv"
	// "time"

	// "github.com/golang-jwt/jwt"
	r "soporte-go/api/routes"
	"soporte-go/core/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type EmpresaHandler struct {
	EmpresaUseCase empresa.EmpresaUseCase
}

func NewEmpresaHandler(e *echo.Echo, u empresa.EmpresaUseCase) {
	handler := &EmpresaHandler{
		EmpresaUseCase: u,
	}

	e.GET("empresa/", handler.GetEmpresaUser)
	e.GET("empresa/areas/", handler.GetAreasFromUser)
	e.GET("empresa/areas-user/", handler.GetAreasUser)

	e.GET("empresa/areas/:areaName/", handler.GetAreaDetail)
	e.GET("empresa/proyecto/:proyectoName/", handler.GetProyectoDetail)

	e.POST("empresa/", handler.StoreEmpresa)
	e.POST("empresa/create-area/", handler.StoreArea)
	e.POST("empresa/add-user-to-area/", handler.AddUserToArea)
	e.GET("empresa/area-change-state/:areaId/:areaState/", handler.AreaChangeState)
	e.GET("empresa/empresa-by-parent-id/", handler.GetEmpresasUser)
	e.POST("empresa/create-proyecto/", handler.CreateProyecto)
	e.GET("empresa/proyectos/:parentId/", handler.GetSubAreas)
	e.GET("empresa/areas-empresa/:empresaId/", handler.GetAreasEmpresa)
	e.GET("empresa/funcionarios-by-area/:areaId/", handler.GetFuncionariosByAreaId)
	e.GET("empresa/clientes-by-area/:areaId/", handler.GetClientesByAreaId)

}

func (u *EmpresaHandler) GetClientesByAreaId(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	areaId, _ := strconv.Atoi(c.Param("areaId"))
	ctx := c.Request().Context()
	// parentId,_:= strconv.Atoi(c.Param("parentId"))
	res, err := u.EmpresaUseCase.GetClientesByAreaId(ctx, areaId)
	if err != nil {
		// logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) GetFuncionariosByAreaId(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	areaId, _ := strconv.Atoi(c.Param("areaId"))
	ctx := c.Request().Context()
	// parentId,_:= strconv.Atoi(c.Param("parentId"))
	res, err := u.EmpresaUseCase.GetFuncionariosByAreaId(ctx, areaId)
	if err != nil {
		// logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) GetAreasEmpresa(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	empresaId, _ := strconv.Atoi(c.Param("empresaId"))
	log.Println(empresaId)
	res, err := u.EmpresaUseCase.GetAreas(ctx, claims.Empresa,claims.UserId,claims.Rol)
	if err != nil {
		// logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) GetSubAreas(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	parentId, _ := strconv.Atoi(c.Param("parentId"))
	res, err := u.EmpresaUseCase.GetProyectos(ctx, parentId)
	if err != nil {
		// logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) CreateProyecto(c echo.Context) (err error) {
	var subArea empresa.Proyecto
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	err = c.Bind(&subArea)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	subArea.CreadorId = claims.UserId
	subArea.EmpresaParentId = claims.Empresa
	// log.Println(subArea.Start)
	// date, err := time.Parse("2006-01-02", subArea.Start)	
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Println(date)
	// s :="15:49:55 2023-04-14T15:49"
	ctx := c.Request().Context()
	err = u.EmpresaUseCase.CreateProyecto(ctx, &subArea)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, subArea)
}

func (u *EmpresaHandler) GetEmpresasUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetEmpresas(ctx, &claims.Empresa)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) AddUserToArea(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var addUserRequestData empresa.AddUserRequestData
	err = c.Bind(&addUserRequestData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	err = u.EmpresaUseCase.AddUserToArea(ctx, &addUserRequestData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	log.Println(err)
	return c.JSON(http.StatusOK, nil)
	// userIds := c.FormValue("ids")
	// ctx := c.Request().Context()
	// err := u.EmpresaUseCase.AddUser/ToArea(ctx)
}

func (u *EmpresaHandler) AreaChangeState(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	a := c.Param("areaId")
	s := c.Param("areaState")
	id, err := strconv.Atoi(a)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	state, err := strconv.Atoi(s)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	log.Println(id)
	log.Println(state)
	ctx := c.Request().Context()
	err = u.EmpresaUseCase.AreaChangeState(ctx, state, id)
	if err != nil {
		return c.JSON(http.StatusNotModified, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, "Ok")

}

func (u *EmpresaHandler) GetAreaDetail(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	n := c.Param("areaName")
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetAreaByName(ctx, n)
	// empresa.Id = casoId
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) GetProyectoDetail(c echo.Context) (err error) {
	// token := c.Request().Header["Authorization"][0]
	// _, err = r.ExtractClaims(token)
	// if err != nil {
		// return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	// }
	n := c.Param("proyectoName")
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetProyectoByName(ctx, n)
	// empresa.Id = casoId
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) StoreArea(c echo.Context) (err error) {
	var area empresa.Area
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	err = c.Bind(&area)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	area.CreadorId = claims.UserId
	area.EmpresaId = claims.Empresa
	ctx := c.Request().Context()
	err = u.EmpresaUseCase.StoreArea(ctx, &area)
	// empresa.Id = casoId
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, area)
}

func (u *EmpresaHandler) GetAreasFromUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetAreasFromUser(ctx, claims.UserId, claims.Empresa, claims.Rol)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) GetAreasUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetAreasUser(ctx, claims.UserId)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) GetEmpresaUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.EmpresaUseCase.GetEmpresa(ctx, claims.UserId, claims.Rol)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (u *EmpresaHandler) StoreEmpresa(c echo.Context) (err error) {
	var empresa empresa.Empresa

	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	// token := c.Request().Header["Authorization"][0]
	// userId,err := r.ExtractClaims(token)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	// }
	err = c.Bind(&empresa)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	empresa.ParentId = &claims.Empresa
	ctx := c.Request().Context()
	err = u.EmpresaUseCase.StoreEmpresa(ctx, &empresa)
	log.Println(err)
	// empresa.Id = casoId
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, empresa)

}
