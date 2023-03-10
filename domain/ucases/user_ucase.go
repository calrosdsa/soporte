package ucases

import (
	"bytes"
	"context"

	// "errors"
	"fmt"
	"log"
	"net/smtp"
	"soporte-go/core/model/user"
	"time"

	// "log"
	"html/template"
)

type userUseCase struct {
	userRepo       user.UserRepository
	contextTimeout time.Duration
	smtpAuth       smtp.Auth
}

func NewUserUseCases(u user.UserRepository, timeout time.Duration, s smtp.Auth) user.UserUseCases {
	return &userUseCase{
		userRepo:       u,
		contextTimeout: timeout,
		smtpAuth:       s,
	}
}

func (a *userUseCase) DeleteInvitation(ctx context.Context, m string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	err = a.userRepo.DeleteInvitation(ctx, m)
	return
}

func (a *userUseCase) ValidateEmail(ctx context.Context, m string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	err = a.userRepo.ValidateEmail(ctx, m)
	return
}

func (a *userUseCase) sendEmail(emails []string, url string) {
	log.Println(emails)
	// go func() {
		t, _ := template.ParseFiles("templates/register-invitation.html")
		var body bytes.Buffer
		headers := "MIME-version: 1.0;\nContent-Type: text/html;"
		body.Write([]byte(fmt.Sprintf("Subject: yourSubject\n%s\n\n", headers)))
		t.Execute(&body, struct {
			URL string
		}{
			URL: url,
		})
		err := smtp.SendMail("smtp.gmail.com:587", a.smtpAuth, "jorgemiranda0180@gmail.com", emails, body.Bytes())
		if err != nil{
			log.Println(err)
		}
		log.Println("No error")
	// }()
}

func (a *userUseCase) ReSendEmail(m []string, url string) {
	a.sendEmail(m, url)
}

func (a *userUseCase) UserRegisterInvitation(ctx context.Context, url string, to user.UserRegistrationRequest, id string) (res []user.UserShortInfo, err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	invitations := make([]user.UserShortInfo, len(to.To))
	for index, value := range to.To {
		t := user.UserShortInfo{
			Nombre:  value,
			Id:      id,
			IsAdmin: to.IsAdmin,
		}
		val, _ := a.userRepo.CreateUserInvitation(ctx, &t)
		// if err != nil{
		// 	return nil,err
		// }
		invitations[index] = val
	}
	a.sendEmail(to.To, url)

	// go func() {
	// 	t,_ := template.ParseFiles("templates/register-invitation.html")
	// 	var body bytes.Buffer
	// 	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	// 	body.Write([]byte(fmt.Sprintf("Subject: yourSubject\n%s\n\n",headers)))
	// 	t.Execute(&body,struct{
	// 		URL  string
	// 		}{
	// 			URL: url,
	// 		})
	// 		smtp.SendMail("smtp.gmail.com:587",a.smtpAuth,"jorgemiranda0180@gmail.com",to.To,body.Bytes())
	// }()
	return invitations, nil
}

func(u *userUseCase) SearchUser(ctx context.Context,id string,q string)(res []user.UserShortInfo,err error){
	ctx,cancel :=  context.WithTimeout(ctx,u.contextTimeout)
	defer cancel()
	res,err = u.userRepo.SearchUser(ctx,id,q)
	return
}

func (u *userUseCase) GetUsersShortIInfo(ctx context.Context, id string) (res []user.UserShortInfo, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	list, err := u.userRepo.GetUsersShortIInfo(ctx, id)
	list2, err := u.userRepo.GetInvitaciones(ctx, id)
	res = append(list, list2...)
	return
}

func (u *userUseCase) GetClientesByArea(ctx context.Context, id int) (res []user.UserArea, err error) {
	res, err = u.userRepo.GetClientesByArea(ctx, id)
	return
}

func (u *userUseCase) GetClientesFiltered(ctx context.Context,f int)(res []user.UserArea,err error) {
	ctx,cancel := context.WithTimeout(ctx,u.contextTimeout)
	defer cancel()
	res,err = u.userRepo.GetClientesFiltered(ctx,f)
	return
}

func (u *userUseCase) GetFuncionarios(ctx context.Context) ([]user.Funcionario, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	list, err := u.userRepo.GetFuncionarios(ctx)
	if err != nil {
		log.Println(err)
	}
	// log.Println(list)
	return list, err
}

func (u *userUseCase) GetClientes(ctx context.Context, id string) ([]user.UserShortInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	list, err := u.userRepo.GetUsersShortIInfo(ctx, id)
	if err != nil {
		log.Println(err)
	}
	// log.Println(list)
	return list, err

}

func (u *userUseCase) UpdateCliente(ctx context.Context, columns []string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepo.UpdateCliente(ctx, columns, values...)
	if err != nil {
		log.Println(err)
		return err
	}
	// query := `update clientes set nombre = &1,apellido = &2,email=&3,celular=&4,telefono=&5,profile_photo=&6)`
	return nil
}

func (u *userUseCase) UpdateFuncionario(ctx context.Context, columns []string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepo.UpdateCliente(ctx, columns, values...)
	if err != nil {
		log.Println(err)
		return err
	}
	// query := `update clientes set nombre = &1,apellido = &2,email=&3,celular=&4,telefono=&5,profile_photo=&6)`
	return nil
}

func (u *userUseCase) GetClienteById(ctx context.Context, id string) (res user.Cliente, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	res, err = u.userRepo.GetClienteById(ctx, id)
	if err != nil {
		return
	}
	return res, err
}

func (u *userUseCase) GetFuncionarioById(ctx context.Context, id string) (res user.Funcionario, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	res, err = u.userRepo.GetFuncionarioById(ctx, id)
	if err != nil {
		return
	}
	return res, err
}
