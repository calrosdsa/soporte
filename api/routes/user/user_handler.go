package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	// errorType "soporte-go/data/model"
	"soporte-go/core/model/user"
	// "soporte-go/core/repository"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	r "soporte-go/api/routes"
	"soporte-go/core/model"
	// "github.com/sirupsen/logrus"
	// validator "gopkg.in/go-playground/validator.v9"
)

type UserHandler struct {
	UserUcase user.UserUseCases
}

func NewUserHandler(e *echo.Echo, us user.UserUseCases) {
	handler := &UserHandler{
		UserUcase: us,
	}
	
	e.POST("user/update-cliente/", handler.UpdateCliente)
	e.GET("user/cliente/:clienteId/", handler.GetClienteById)
	e.GET("user/clientes/", handler.GetClientes)
	e.POST("user/update-funcionario/", handler.UpdateFuncionario)
	e.GET("user/funcionario/:funcionarioId/", handler.GetClienteById)
	e.GET("user/funcionarios/", handler.GetClientes)
	e.GET("user/clientes-area/:areaId/", handler.GetClientesByArea)
	e.POST("user/register-invitation/", handler.UserRegisterInvitation)
	e.GET("user/user-list/", handler.GetUserList)
	e.GET("user/validate-email/:email/", handler.ValidateEmail)
	e.GET("user/resend-email/", handler.ResendEmail)
	e.GET("user/cancel-invitation/", handler.CancelInvitation)
	e.GET("user/search/", handler.SearchUser)
	e.GET("user/add-user-list/:areaId/", handler.GetUserFiltered)
	e.GET("user/users-empresa/:emId/", handler.GetUsersbyEmpresaId)
	e.GET("user/users-empresa-by-rol/:emId/:rol/", handler.GetUsersEmpresaByRol)

	e.GET("user/", handler.GetProfile)

}

func (u *UserHandler) GetProfile(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.UserUcase.GetUserById(ctx, claims.UserId,claims.Rol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}


func (u *UserHandler) GetUsersEmpresaByRol(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	rol, _ := strconv.Atoi(c.Param("rol"))
	ctx := c.Request().Context()
	emId, _ := strconv.Atoi(c.Param("emId"))
	res, err := u.UserUcase.GetUsersEmpresaByRol(ctx, emId, rol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *UserHandler) GetUsersbyEmpresaId(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	emId, _ := strconv.Atoi(c.Param("emId"))
	res, err := u.UserUcase.GetUsersbyEmpresaId(ctx, emId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *UserHandler) GetUserFiltered(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	areaId := c.Param("areaId")
	id, _ := strconv.Atoi(areaId)
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := a.UserUcase.GetUserAddList(ctx, id, claims.Rol, claims.UserId)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *UserHandler) SearchUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	q := c.QueryParam("search")
	ctx := c.Request().Context()
	res, err := a.UserUcase.SearchUser(ctx, claims.UserId, q)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *UserHandler) CancelInvitation(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	m := c.QueryParam("email")
	log.Println(m)
	ctx := c.Request().Context()
	err = a.UserUcase.DeleteInvitation(ctx, m)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (a *UserHandler) ResendEmail(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	m := c.QueryParam("email")
	ad := c.QueryParam("isAdmin")
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var rol int

	if claims.Rol == int(model.RoleClienteAdmin) {
		if ad == "true" {
			rol = 2
		} else {
			rol = 0
		}
	} else if claims.Rol == int(model.RoleFuncionarioAdmin) {
		if ad == "true" {
			rol = 3
		} else {
			rol = 1
		}
	} else {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: errors.New("no rol presente en jwt token").Error()})
	}
	tokenInvitation, err := r.GenerateInvitationJWT(claims.UserId, rol, claims.Empresa, m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	url := fmt.Sprintf("http://localhost:3000/auth/registro?auth=%s", tokenInvitation)
	a.UserUcase.ReSendEmail([]string{m}, url)
	return c.JSON(http.StatusOK, nil)

}

func (a *UserHandler) ValidateEmail(c echo.Context) (err error) {
	m := c.Param("email")
	ctx := c.Request().Context()
	err = a.UserUcase.ValidateEmail(ctx, m)
	if err != nil {
		return c.JSON(http.StatusConflict, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, "")
}

func (a *UserHandler) GetUserList(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := a.UserUcase.GetUsersShortIInfo(ctx, claims.UserId, claims.Rol, claims.Empresa)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *UserHandler) UserRegisterInvitation(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var to user.UserRegistrationRequest
	err = c.Bind(&to)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	var rol int
	log.Println(to.IsAdmin)
	if to.EmpresaId != 0 {
		if to.IsAdmin {
			rol = 2
		} else {
			rol = 0
		}
		claims.Empresa = to.EmpresaId
	} else {
		if claims.Rol == int(model.RoleClienteAdmin) {
			if to.IsAdmin {
				rol = 2
			} else {
				rol = 0
			}
			} else if claims.Rol == int(model.RoleFuncionarioAdmin) {
			if to.IsAdmin {
				rol = 3
			} else {
				rol = 1
			}
		} else {
			return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: errors.New("no rol presente en jwt token").Error()})
		}
		// log.Println(rol)
	}
	
	log.Println(claims.Empresa)
	ctx := c.Request().Context()
	res, err := a.UserUcase.UserRegisterInvitation(ctx, &to, claims.UserId, rol, claims.Empresa)
	// log.Println(url)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)

}

func (u *UserHandler) GetClientesByArea(c echo.Context) (err error) {
	idS := c.Param("areaId")
	id, err := strconv.Atoi(idS)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.UserUcase.GetClientesByArea(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (u *UserHandler) GetFuncionarioById(c echo.Context) (err error) {
	id := c.Param("funcionarioId")
	ctx := c.Request().Context()
	res, err := u.UserUcase.GetFuncionarioById(ctx, id)
	if err != nil {
		logrus.Error(err)
	}
	mar, _ := json.Marshal(&res)
	log.Println(string(mar))
	return c.JSON(http.StatusOK, res)
}

func (u *UserHandler) UpdateFuncionario(c echo.Context) (err error) {
	// var cliente user.Cliente
	id := c.QueryParam("id")
	var json map[string]interface{} = map[string]interface{}{}
	err = c.Bind(&json)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	var columns []string
	var values []interface{}
	for k, v := range json {
		columns = append(columns, k)
		values = append(values, v)
		log.Println(reflect.TypeOf(v))
	}
	values = append(values, id)
	// log.Println(fmt.Sprintf("%v", json))

	ctx := c.Request().Context()
	err = u.UserUcase.UpdateFuncionario(ctx, columns, values...)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, "okk")
}

func (u *UserHandler) GetClientes(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res, err := u.UserUcase.GetClientes(ctx, claims.UserId, claims.Rol)
	if err != nil {
		logrus.Error(err)
	}
	return c.JSON(http.StatusOK, res)
}

func (u *UserHandler) GetClienteById(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}

	id := c.Param("clienteId")
	ctx := c.Request().Context()
	res, err := u.UserUcase.GetUserById(ctx, id, claims.Rol)
	if err != nil {
		logrus.Error(err)
	}
	mar, _ := json.Marshal(&res)
	log.Println(string(mar))
	return c.JSON(http.StatusOK, res)
}

func (u *UserHandler) UpdateCliente(c echo.Context) (err error) {
	// var cliente user.Cliente
	clientId := c.QueryParam("id")
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}

	var json map[string]interface{} = map[string]interface{}{}
	err = c.Bind(&json)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	log.Println(len(json))
	var columns []string
	var values []interface{}
	for k, v := range json {

		columns = append(columns, k)
		values = append(values, v)
		log.Println(v)
		log.Println(k)
	}
	log.Println(clientId)
	values = append(values, clientId)
	// log.Println(fmt.Sprintf("%v", json))

	ctx := c.Request().Context()
	err = u.UserUcase.UpdateCliente(ctx, columns, values...)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, "okk")
}
