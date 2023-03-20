package routes

import (
	"net/http"
	"soporte-go/core/model/media"

	"github.com/labstack/echo/v4"

	r "soporte-go/api/routes"
	"soporte-go/core/model"
)

type MediaHandler struct {
	MUseCase media.MediaUseCase
}

func NewMediaHandler(e *echo.Echo, mu media.MediaUseCase) {
	handler := &MediaHandler{
		MUseCase: mu,
	}

	e.POST("media/upload-file-caso/", handler.UploadFileCaso)
	e.GET("media/caso-files/:id/",handler.GetFileCasos)
}

func (m *MediaHandler) GetFileCasos(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	id := c.Param("id")
	ctx := c.Request().Context()
	res,err := m.MUseCase.GetFileCasos(ctx,id)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (m *MediaHandler) UploadFileCaso(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	casoId := c.FormValue("casoId")
	descripcion := c.FormValue("descripcion")
	extension := c.FormValue("extension")
	ctx := c.Request().Context()
	fileCaso, err := m.MUseCase.UploadFileCaso(ctx, file, casoId, descripcion, extension)
	if err != nil {
		return c.JSON(model.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, fileCaso)

}
