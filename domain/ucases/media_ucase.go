package ucases

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"soporte-go/core/model/media"
	"soporte-go/util"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nickalie/go-webpbin"
)

// const

type mediaUseCase struct {
	mediaRepo      media.MediaRepository
	contextTimeout time.Duration
	sess           *session.Session
}

func NewMediaUseCase(mu media.MediaRepository, timeout time.Duration, sess *session.Session) media.MediaUseCase {
	return &mediaUseCase{
		mediaRepo:      mu,
		contextTimeout: timeout,
		sess:           sess,
	}
}
func (mu *mediaUseCase) UploadImageWithoutCtx(wg *sync.WaitGroup,file *multipart.FileHeader,filename string,urls *[]string) {
	src, err := file.Open()
	if err != nil {
		return
	}
	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return
	}
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err)
	}

	err = webpbin.NewCWebP().
		Quality(40).
		InputFile(dst.Name()).
		OutputFile(filename).
		Run()
	if err != nil {
		log.Println(err)
	}
	fileWebp, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	url, err := util.UplaodObjectWebpWithoutCxt(fileWebp, "teclu-soporte", mu.sess)
	if err != nil {
		log.Println(err)
	}
	*urls = append(*urls, url)
	
	defer func() {
		src.Close()
		if err := dst.Close(); err != nil {
			log.Println(err)
		}
		err := os.Remove(dst.Name())
		if err != nil {
			log.Println(err)
		}
		if err := fileWebp.Close(); err != nil {
			log.Println(err)
		}
		err1 := os.Remove(filename)
		if err1 != nil {
			log.Println(err1)
		}
		wg.Done()
	}()
	return	
}

func (mu *mediaUseCase) UploadImage(ctx context.Context, file *multipart.FileHeader, filename string) (url string, err error) {
	ctx, cancel := context.WithTimeout(ctx, mu.contextTimeout)
	defer cancel()
	src, err := file.Open()
	if err != nil {
		return
	}
	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return
	}
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err)
	}

	err = webpbin.NewCWebP().
		Quality(40).
		InputFile(dst.Name()).
		OutputFile(filename).
		Run()
	if err != nil {
		log.Println(err)
	}
	fileWebp, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	url, err = util.UplaodObjectWebp(ctx,fileWebp, "teclu-soporte", mu.sess)
	
	defer func() {
		src.Close()
		if err := dst.Close(); err != nil {
			log.Println(err)
		}
		err := os.Remove(dst.Name())
		if err != nil {
			log.Println(err)
		}
		if err := fileWebp.Close(); err != nil {
			log.Println(err)
		}
		err1 := os.Remove(filename)
		if err1 != nil {
			log.Println(err1)
		}
	}()
	return
}

func (mu *mediaUseCase) UploadFileCaso(ctx context.Context, file *multipart.FileHeader, id string, descripcion string, ext string) (res media.CasoFile, err error) {
	ctx, cancel := context.WithTimeout(ctx, mu.contextTimeout)
	defer cancel()
	url, err := util.UplaodObject(file, "teclu-soporte", mu.sess)
	if err != nil {
		return
	}
	res, err = mu.mediaRepo.UploadFileCaso(ctx, url, id, descripcion, ext)
	return
}

func (mu *mediaUseCase) GetFileCasos(ctx context.Context, id string) (res []media.CasoFile, err error) {
	ctx, cancel := context.WithTimeout(ctx, mu.contextTimeout)
	defer cancel()
	res, err = mu.mediaRepo.GetFileCasos(ctx, id)
	return
}
