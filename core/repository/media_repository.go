package repository

import (
	"context"
	"database/sql"
	"soporte-go/core/model/media"
	"time"

	"github.com/sirupsen/logrus"
)

type mediaRepository struct {
	Conn *sql.DB
	Context context.Context
}


func NewMediaRepository(conn *sql.DB,ctx context.Context) media.MediaRepository {
	return &mediaRepository{
		Conn: conn,
		Context: ctx,
	}
}

func (p *mediaRepository) UploadFileCaso(ctx context.Context,url string,id string,descripcion string,ext string)(res media.CasoFile,err error){
	query := `insert into recursos (file_url,ext,descripcion,caso_id,created_on) values ($1,$2,$3,$4,$5)
	returning (id,file_url,ext,descripcion,caso_id,created_on);`
	t := media.CasoFile{}
	err =  p.Conn.QueryRowContext(ctx,query,url,ext,descripcion,id,time.Now()).Scan(&t)
	if err != nil{
		return
	}
	return t,nil
}

func (p *mediaRepository) GetFileCasos(ctx context.Context,id string)(res []media.CasoFile,err error){
	query := `select * from recursos where caso_id = $1;`
	res,err = p.fetchFileCasos(ctx,query,id)
	return
}

func (p *mediaRepository) fetchFileCasos(ctx context.Context,query string,args ...interface{})(result []media.CasoFile,err error){
	rows, err := p.Conn.QueryContext(p.Context, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]media.CasoFile, 0)
	for rows.Next() {
		t := media.CasoFile{}
		err = rows.Scan(
			&t.Id,
			&t.FileUrl,
			&t.Extension,
			&t.Descripcion,
			&t.CasoId,
			&t.CreatedOn,
		)
		result = append(result, t)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return result, nil
}