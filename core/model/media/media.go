package media

import (
	"context"
	"mime/multipart"
	"sync"
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
	UploadImage(ctx context.Context,file *multipart.FileHeader,filename string)(url string,err error)
	UploadImageWithoutCtx(wg *sync.WaitGroup,file *multipart.FileHeader,filename string,urls *[]string)
}

type MediaRepository interface {
	UploadFileCaso(ctx context.Context,url string,id string,descripcion string,ext string) (CasoFile,error)
	GetFileCasos(ctx context.Context,id string)([]CasoFile,error)
}