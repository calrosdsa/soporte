package ucases

import (
	"context"
	"mime/multipart"
	"soporte-go/core/model/media"
	"soporte-go/util"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
)

type mediaUseCase struct {
	mediaRepo      media.MediaRepository
	contextTimeout time.Duration
	sess           *session.Session
}

func NewMediaUseCase (mu media.MediaRepository ,timeout time.Duration, sess *session.Session) media.MediaUseCase{
	return &mediaUseCase{
		mediaRepo: mu,
		contextTimeout: timeout,
		sess: sess,
	}
}

func (mu *mediaUseCase) UploadFileCaso(ctx context.Context,file *multipart.FileHeader,id string,descripcion string,ext string)(res media.CasoFile,err error){
	ctx, cancel := context.WithTimeout(ctx,mu.contextTimeout)
	defer cancel()
	url,err := util.UplaodObject(file, "teclu-soporte", mu.sess)
	if err != nil {
		return 
	}
	res,err = mu.mediaRepo.UploadFileCaso(ctx,url,id,descripcion,ext)
	return
}

func (mu *mediaUseCase) GetFileCasos(ctx context.Context,id string)(res []media.CasoFile,err error){
	ctx,cancel := context.WithTimeout(ctx,mu.contextTimeout)
	defer cancel()
	res,err = mu.mediaRepo.GetFileCasos(ctx,id)
	return
}