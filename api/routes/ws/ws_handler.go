package routes

import (
	"context"
	"log"
	"net/http"
	"soporte-go/core/model/ws"
	_ws "soporte-go/ws"
	"time"

	_routes "soporte-go/api/routes"
	"soporte-go/core/model"

	"github.com/labstack/echo/v4"
)

type WsHandler struct {
	re ws.WsRepository
}

func NewWsHandler(e *echo.Echo, re ws.WsRepository) {
	
	go _ws.H.Run(re)
	handler := &WsHandler{
		re: re,
	}
	e.GET("/ws/:casoId", handler.ChatCaso)
	e.GET("/ws/messages/:casoId",handler.GetMessages)
}

func (ws *WsHandler) GetMessages(c echo.Context)(err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = _routes.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	casoId := c.Param("casoId")
	ctx := c.Request().Context()
	timeoutContext := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(ctx,timeoutContext)
	defer cancel()
	res,err := ws.re.GetMessages(ctx,casoId)
	if err != nil {
		// logrus.Error(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (ws *WsHandler) ChatCaso(c echo.Context) (err error) {
	casoId := c.Param("casoId")
	log.Println("casoId..", casoId)
	_ws.ServeWs(c.Response(), c.Request(), casoId)
	return nil
}
