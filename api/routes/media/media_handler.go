package routes

import (
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"soporte-go/core/model/media"
	"sync"

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
	e.GET("media/caso-files/:id/", handler.GetFileCasos)
	e.POST("media/upload-image/", handler.UploadImage)
	e.POST("media/upload-multiple-images/", handler.UploadMultipleImages)
}

func (m *MediaHandler) UploadMultipleImages(c echo.Context) (err error) {
	form, _ := c.MultipartForm()
	files := form.File["images"]
	var wg sync.WaitGroup
	wg.Add(len(files))
	urls := []string{}
	for _, file := range files {
		// log.Println(filepath.Ext(file.Filename))
		go func(file *multipart.FileHeader) {
			filename := file.Filename[0:len(file.Filename)-len(filepath.Ext(file.Filename))] + ".webp"
			log.Println(filename)
			m.MUseCase.UploadImageWithoutCtx(&wg, file, filename,&urls)
		}(file)
	}
	wg.Wait()
	log.Println(len(files))
	return c.JSON(200, urls)
}

func (m *MediaHandler) UploadImage(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	file, err := c.FormFile("file")
	if err != nil {
		// log.Println(err)
		return c.JSON(http.StatusNotFound, model.ResponseError{Message: err.Error()})
	}
	filename := c.FormValue("filename")

	ctx := c.Request().Context()
	url, err := m.MUseCase.UploadImage(ctx, file, filename)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, url)
}

func (m *MediaHandler) GetFileCasos(c echo.Context) (err error) {
	token := c.Request().Header["Authorization"][0]
	_, err = r.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ResponseError{Message: err.Error()})
	}
	id := c.Param("id")
	ctx := c.Request().Context()
	res, err := m.MUseCase.GetFileCasos(ctx, id)
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
