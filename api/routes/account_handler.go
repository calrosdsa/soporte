package routes

import (
	"errors"
	
	"log"
	"net/http"
	

	// "strconv"
	errorType "soporte-go/core/model"
	"soporte-go/core/model/account"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	validator "gopkg.in/go-playground/validator.v9"
)

type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type AccountHandler struct {
	AUsecase account.AccountUseCase
}

type ValidData interface {
	*account.RegisterForm | *account.LoginRequest
}

func NewAccountHandler(e *echo.Echo, us account.AccountUseCase) {
	handler := &AccountHandler{
		AUsecase: us,
	}
	e.POST("account/register-cliente/", handler.RegisterCliente)
	e.POST("account/register-funcionario/",handler.RegisterFuncionario)
	e.POST("account/register-cliente-admin/",handler.RegisterClienteAdministrador)
	e.POST("account/login/", handler.Login)
	e.GET("account/delete-user/",handler.DeleteUser)
}

func(a *AccountHandler) DeleteUser(c echo.Context)( err error){
	token := c.Request().Header["Authorization"][0]
	_, err = ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	id := c.QueryParam("id")
	ctx := c.Request().Context()
	err = a.AUsecase.DeleteUser(ctx,id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (a *AccountHandler) RegisterClienteAdministrador(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	if claims.Rol != 3 && claims.Rol !=2 {
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: errors.New("No tiene los permisos para registrar").Error()})
	}
	var user account.RegisterForm
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	
	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user.Rol = 2
	user.EmpresaId = claims.Empresa
	// token,err := jwt.GenerateJWT(user.)
	// log.Println(user)
	ctx := c.Request().Context()
	res, err := a.AUsecase.RegisterCliente(ctx, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, res)
}

func (a *AccountHandler) Login(c echo.Context) (err error) {
	var loginRequest account.LoginRequest
	err = c.Bind(&loginRequest)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	var ok bool
	if ok, err = isRequestValid(&loginRequest); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	res, err := a.AUsecase.Login(ctx, &loginRequest)
	if err != nil {
		logrus.Error(err)
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	token, err := GenerateJWT(res.ClientId, res.Rol,res.EmpresaId)
	if err != nil {
		log.Println(err)
	}
	response := account.AuthenticationResponse{
		Token: token,
		User:  res,
	}

	return c.JSON(http.StatusOK, response)
}


func (a *AccountHandler) RegisterFuncionario(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claims, err := ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	if claims.Rol != 3 {
		return c.JSON(http.StatusUnauthorized, ResponseError{
			Message: errors.New("No tiene los permisos para registrar").Error()})
	}
	var user account.RegisterForm
	user.Rol = 1
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// token,err := jwt.GenerateJWT(user.)
	// log.Println(user)
	ctx := c.Request().Context()
	res, err := a.AUsecase.RegisterFuncionario(ctx, &user, claims.UserId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, res)
}

func (a *AccountHandler) RegisterCliente(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	claimsInvitation, err := ExtractClaimsInvitation(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	// if claims.Rol == 0 {
	// 	return c.JSON(http.StatusUnauthorized, ResponseError{
	// 		Message: errors.New("No tiene los permisos para registrar clientes").Error()})
	// }

	var userForm account.RegisterForm
	err = c.Bind(&userForm)
	userForm.Rol = claimsInvitation.Rol
	userForm.EmpresaId = claimsInvitation.EmpresaId
	userForm.SuperiorId = claimsInvitation.Id
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&userForm); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// log.Println(userForm)
	ctx := c.Request().Context()
	res, err := a.AUsecase.RegisterCliente(ctx, &userForm)
	res.Rol = claimsInvitation.Rol
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}
	token,err = GenerateJWT(*res.UserId,userForm.Rol,userForm.EmpresaId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	response := account.RegisterAuthResponse{
		Access: token,
		User: res,
	}
	return c.JSON(http.StatusCreated, response)
}

func isRequestValid[T ValidData](m T) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {

	case errorType.ErrInternalServerError:
		return http.StatusInternalServerError
	case errorType.ErrNotFound:
		return http.StatusNotFound
	case errorType.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
