package media

import (
	"context"
	"mime/multipart"
	"time"
)
type CasoFile struct{
	Id int `json:"id"`
	FileUrl string `json:"file_url"`
	Extension string `json:"extension"`
	Descripcion string `json:"descripcion"`
	CasoId string `json:"caso_id"`
	CreatedOn time.Time `json:"created_on"`
}


type MediaUseCase interface{
	UploadFileCaso(ctx context.Context,file *multipart.FileHeader,id string,descripcion string,ext string) (CasoFile,error)
	GetFileCasos(ctx context.Context,id string)([]CasoFile,error)
}

type MediaRepository interface {
	UploadFileCaso(ctx context.Context,url string,id string,descripcion string,ext string) (CasoFile,error)
	GetFileCasos(ctx context.Context,id string)([]CasoFile,error)
}