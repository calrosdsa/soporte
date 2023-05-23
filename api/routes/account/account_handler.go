package account

import (
	// "errors"

	"log"
	"net/http"

	// "strconv"
	"soporte-go/core/model"
	"soporte-go/core/model/account"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	r "soporte-go/api/routes"

	validator "gopkg.in/go-playground/validator.v9"
)

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
	e.POST("account/register/",handler.RegisterUser)
	 
	// e.POST("account/register-cliente/", handler.RegisterCliente)
	// e.POST("account/register-funcionario/", handler.RegisterFuncionario)
	// e.POST("account/register-cliente-admin/", handler.RegisterClienteAdministrador)
	e.POST("account/login/", handler.Login)
	e.POST("account/update-password/", handler.UpdatePassword)
	e.GET("account/delete-user/", handler.DeleteUser)

	// e.GET("account/send-reset-password/:email/",)
}

// func (a *AccountHandler) SendResetPassoword(c echo.Context)(err error){
// 	email := c.Param("email")
// 	token, err := r.GenerateEmailJWT(email)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	// a.AUsecase.
// }

func (a *AccountHandler) UpdatePassword(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var data account.PasswordUpdate
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	
	ctx := c.Request().Context()
	err = a.AUsecase.UpdatePassword(ctx,data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, "Se ha actualizado con exito")
}


func (a *AccountHandler) RegisterUser(c echo.Context) (err error) {
	tokenInvitation := c.Request().Header["Authorization"][0]
	claims, err := r.ExtractClaimsInvitation(tokenInvitation)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	var user account.RegisterForm
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	user.Email = claims.Email
	user.Rol = claims.Rol
	user.EmpresaId = claims.EmpresaId
	user.SuperiorId = &claims.Id

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// token,err := jwt.r.GenerateJWT(user.)
	// log.Println(user)
	ctx := c.Request().Context()
	res, err := a.AUsecase.RegisterUser(ctx, &user)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	token, err := r.GenerateJWT(res.Id, res.Rol, res.EmpresaId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	response := account.AuthenticationResponse{
		Token: token,
		User:  res,
	}
	return c.JSON(http.StatusCreated, response)
}

func (a *AccountHandler) DeleteUser(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	id := c.QueryParam("id")
	ctx := c.Request().Context()
	err = a.AUsecase.DeleteUser(ctx, id)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
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
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	token, err := r.GenerateJWT(res.Id, res.Rol, res.EmpresaId)
	if err != nil {
		log.Println(err)
	}
	response := account.AuthenticationResponse{
		Token: token,
		User:  res,
	}

	return c.JSON(http.StatusOK, response)
}


func isRequestValid[T ValidData](m T) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
