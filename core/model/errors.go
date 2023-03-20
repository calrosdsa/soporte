package model

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")
	
	ErrConflictEmail = errors.New("este correo electrónico ya ha sido registrado por otro usuario")

	ErrConflictUsername = errors.New("lo siento, ese nombre de usuario ya ha sido tomado")

)


type ResponseError struct {
	Message string `json:"message"`
}



func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {

	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	case ErrConflictEmail:
		return http.StatusConflict
	case ErrConflictUsername:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
