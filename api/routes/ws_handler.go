package routes

import (
	"log"
	"soporte-go/core/model/ws"
	_ws "soporte-go/ws"

	"github.com/labstack/echo/v4"
)

type WsHandler struct {
	// re ws.WsRepository
}

func NewWsHandler(e *echo.Echo, re ws.WsRepository) {
	
	go _ws.H.Run(re)
	handler := &WsHandler{
		// re: re,
	}
	e.GET("/ws/:casoId", handler.ChatCaso)
}

func (ws *WsHandler) ChatCaso(c echo.Context) (err error) {
	casoId := c.Param("casoId")
	log.Println("casoId..", casoId)
	_ws.ServeWs(c.Response(), c.Request(), casoId)
	return nil
}
