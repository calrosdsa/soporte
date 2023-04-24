package caso

import (
	"fmt"
	"log"
	"net/http"

	// "reflect"
	model "soporte-go/core/model"
	"soporte-go/core/model/caso"
	"strconv"

	// "github.com/golang-jwt/jwt"
	_routes "soporte-go/api/routes"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type CasoHandler struct {
	CasoUseCase caso.CasoUseCase
}

func NewCasoHandler(e *echo.Echo, uc caso.CasoUseCase) {
	handler := &CasoHandler{
		CasoUseCase: uc,
	}

	e.POST("/caso/", handler.CreateCaso)
	e.GET("/caso/:casoId/:rol/", handler.GetCaso)
	e.GET("/casos", handler.GetCasosUser)
	e.GET("/caso/from-user-caso/", handler.GetCasosFromUserCaso)
	e.GET("/casos-all/", handler.GetAllCasosUser)
	e.POST("/caso/caso-update/", handler.UpdateCaso)

	e.GET("/caso/asignar-funcionario/:casoId/:idFuncionario/", handler.AsignarFuncionario)
	e.POST("/caso/finalizar-caso/", handler.FinalizarCaso)

	e.POST("/caso/reporte-casos-xlsx/", handler.GetReporteCasosXlsx)
	e.POST("/caso/reporte-caso-html/", handler.GetReporteHtml)

	e.POST("/caso/add-user-caso/", handler.AsignarFuncionarioSoporte)
	e.POST("/caso/usuarios-caso/:casoId/", handler.GetUsuariosCaso)
	// e.GET("/caso/reporte-casos-pdf/",handler.GetReporteCasosPdf)
}

type User struct {
	Name  string
	Age   int
	Email string
}

func (u *CasoHandler) GetReporteHtml(c echo.Context) (err error) {
	// token := c.Request().Header["Authorization"][0]
	// _, err = _routes.ExtractClaims(token)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	// }
	// id := c.Param("casoId")
	log.Println("GET REPORTE HTML")
	var data caso.Caso
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	// bytes,err := u.CasoUseCase.GetReporteCasos(ctx,model.HTML,&data)
	// if err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})

	// }
	buf,err :=u.CasoUseCase.GetReporteCaso(ctx,model.HTML,data)
	if err != nil {
		log.Println(err)
	}
	return c.Blob(http.StatusOK, "reporte.html", buf.Bytes())
}

func (u *CasoHandler) GetUsuariosCaso(c echo.Context) (err error) {
	// token := c.Request().Header["Authorization"][0]
	// _, err = _routes.ExtractClaims(token)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	// }
	id := c.Param("casoId")
	ctx := c.Request().Context()
	res, err := u.CasoUseCase.GetUsuariosCaso(ctx, id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (u *CasoHandler) AsignarFuncionarioSoporte(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var data caso.UserCaso
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	// log.Println(data.Detail)
	ctx := c.Request().Context()
	err = u.CasoUseCase.AsignarFuncionarioSoporte(ctx, &data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "Ok")
}

// func (u *CasoHandler) GetReporteCasosPdf(c echo.Context)(err error) {
// 	ctx := c.Request().Context()
// 	bytes,err := u.CasoUseCase.GetReporteCasos(ctx,model.PDF)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return c.Blob(http.StatusOK,"reporte.pdf",bytes.Bytes())
// }

func (u *CasoHandler) GetReporteCasosXlsx(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var data caso.CasoReporteOptions
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	bytes, err := u.CasoUseCase.GetReporteCasos(ctx, model.XLSX, &data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})

	}
	return c.Blob(http.StatusOK, "reporte.xlsx", bytes.Bytes())
}

func (u *CasoHandler) UpdateCaso(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var data caso.Caso
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	err = u.CasoUseCase.UpdateCaso(ctx, &data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})

	}
	return c.JSON(http.StatusOK, data)
}

func (u *CasoHandler) FinalizarCaso(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var data caso.FinalizacionDetail
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	// log.Println(data.Detail)
	ctx := c.Request().Context()
	err = u.CasoUseCase.FinalizarCaso(ctx, &data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "Ok")
}

func (u *CasoHandler) AsignarFuncionario(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	id := c.Param("casoId")
	idF := c.Param("idFuncionario")
	ctx := c.Request().Context()
	err = u.CasoUseCase.AsignarFuncionario(ctx, id, idF)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, "Ok")
}

func (u *CasoHandler) GetAllCasosUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var casoQuery caso.CasoQuery
	u.CasoQueries(c, &casoQuery)
	ctx := c.Request().Context()
	res, size, err := u.CasoUseCase.GetAllCasosUser(ctx, claims.UserId, &casoQuery, claims.Rol)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	response := caso.CasosResponse{
		Casos:   res,
		Size:    size,
		Current: casoQuery.Page,
	}
	return c.JSON(http.StatusOK, response)
}

func (u *CasoHandler) GetCasosFromUserCaso(c echo.Context) (err error) {
	// page := c.QueryParam("page")
	token := c.Request().Header["Authorization"][0]
	claims, err := _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var casoQuery caso.CasoQuery
	u.CasoQueries(c, &casoQuery)
	ctx := c.Request().Context()
	res, err := u.CasoUseCase.GetCasosFromUserCaso(ctx, claims.UserId, &casoQuery)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	response := caso.CasosResponse{
		Casos:   res,
		Size:    10,
		Current: casoQuery.Page,
	}
	return c.JSON(http.StatusOK, response)
}

func (u *CasoHandler) GetCasosUser(c echo.Context) (err error) {
	// page := c.QueryParam("page")
	token := c.Request().Header["Authorization"][0]
	claims, err := _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var casoQuery caso.CasoQuery
	u.CasoQueries(c, &casoQuery)
	ctx := c.Request().Context()
	res, size, err := u.CasoUseCase.GetCasosUser(ctx, claims.UserId, &casoQuery, claims.Rol)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	response := caso.CasosResponse{
		Casos:   res,
		Size:    size,
		Current: casoQuery.Page,
	}
	return c.JSON(http.StatusOK, response)
}

func (u *CasoHandler) GetCaso(c echo.Context) (err error) {
	id := c.Param("casoId")
	rol := c.Param("rol")
	token := c.Request().Header["Authorization"][0]
	_, err = _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	r, err := strconv.Atoi(rol)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	// token := c.Request().Header["Authorization"][0]
	// userId,err := ExtractClaims(token)
	ctx := c.Request().Context()
	res, err := u.CasoUseCase.GetCaso(ctx, id, r)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *CasoHandler) CreateCaso(c echo.Context) (err error) {
	var caso caso.Caso
	token := c.Request().Header["Authorization"][0]
	claims, err := _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	err = c.Bind(&caso)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	err = u.CasoUseCase.CreateCaso(ctx, &caso, claims.UserId, claims.Empresa, claims.Rol)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, caso)
}

func (u *CasoHandler) CasoQueries(c echo.Context, q *caso.CasoQuery) {
	updated, _ := strconv.Atoi(c.QueryParam("updated"))
	created, _ := strconv.Atoi(c.QueryParam("created"))
	proyecto := c.QueryParam("proyecto")
	key := c.QueryParam("key")
	q.Page, _ = strconv.Atoi(c.QueryParam("page"))
	log.Println("KEYYY", key)

	switch key {
	case "undefined":
		q.Key = ""
	case "":
		q.Key = ""
	default:
		q.Key = fmt.Sprintf(`and c.key = '%s'`, key)
	}

	log.Println(proyecto)
	switch proyecto {
	case "undefined":
		q.Proyecto = ""
	case "0":
		q.Proyecto = ""
	case "":
		q.Proyecto = ""
	default:
		q.Proyecto = fmt.Sprintf(`and c.area = %s`, proyecto)
	}

	switch model.Order(created) {
	case model.ASC:
		q.Order = "order by created_on ASC"
	case model.DESC:
		q.Order = "order by created_on DESC"
	}
	log.Println(q.Order)
	switch model.Order(updated) {
	case model.ASC:
		q.Order = "order by updated_on ASC"
	case model.DESC:
		q.Order = "order by updated_on DESC"
	}

}
