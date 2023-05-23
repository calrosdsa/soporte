package util

import "context"

type UtilRepository interface {
	GetEmailFromUserCasos(ctx context.Context, id int,userId string) (res NotificationCasoEmail, err error)
}

type Mails struct {
	Mail string
}

type NotificationCasoEmail struct {
	Mails []Mails
	Area string 
	Entidad string
	UsuarioName string
	UsuarioApellido string
}