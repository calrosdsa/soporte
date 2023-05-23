package ucases

import (
	"bytes"
	"context"
	"html/template"
	"log"

	"soporte-go/core/model"
	"soporte-go/core/model/account"
	"soporte-go/core/model/user"
	"time"

	"gopkg.in/gomail.v2"
)

type accountUseCase struct {
	accountRepo    account.AccountRepository
	contextTimeout time.Duration
	util           model.Util
	gomail         *gomail.Dialer
}

func NewAccountUseCase(a account.AccountRepository, timeout time.Duration, util model.Util, gomail *gomail.Dialer) account.AccountUseCase {
	return &accountUseCase{
		accountRepo:    a,
		contextTimeout: timeout,
		util:           util,
		gomail:         gomail,
	}
}

func (a *accountUseCase) UpdatePassword(ctx context.Context, d account.PasswordUpdate) (err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	err = a.accountRepo.UpdatePassword(ctx, d)
	return
}

func (a *accountUseCase) RegisterUser(ctx context.Context, form *account.RegisterForm) (res user.UserAuth, err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	if a.util.IsClienteRol(form.Rol) {
		res, err = a.accountRepo.RegisterCliente(ctx, form)
	} else if a.util.IsFuncionarioRol(form.Rol) {
		res, err = a.accountRepo.RegisterFuncionario(ctx, form)
	}
	return
}

func (a *accountUseCase) DeleteUser(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	err = a.accountRepo.DeleteUser(ctx, id)
	return
}

func (a *accountUseCase) Login(ctx context.Context, loginRequest *account.LoginRequest) (res user.UserAuth, err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	res, err = a.accountRepo.Login(ctx, loginRequest)
	if err != nil {
		return
	}
	return
}

func (a *accountUseCase) SendEmailResetPassword(email string,url string) {
	// log.Println(emails)
	// mails := []string{"diegoarmando12ab34cd@gmail.com"}
	go func() {
		t, _ := template.ParseFiles("templates/register-invitation.html")
		var body bytes.Buffer
		t.Execute(&body, struct {
			URL string
		}{
			URL: url,
		})
		m := gomail.NewMessage()
		m.SetHeader("From", "jmiranda@teclu.com")
		m.SetHeader("To", email)
		// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "Restablece tu contraseña")
		m.SetBody("text/html", body.String())
		// m.Attach("/home/Alex/lolcat.jpg")

		// d := gomail.NewDialer("mail.teclu.com", 25, "jmiranda@teclu.com", "jmiranda2022")
		if err := a.gomail.DialAndSend(m); err != nil {
			log.Println("Error sending email", err)
		}

		// // headers := "MIME-version: 1.0;\nContent-Type: text/html;"
		// mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		// body.Write([]byte(fmt.Sprintf("Subject: yourSubject\n%s\n\n", mimeHeaders)))
		// // msg :=[]byte("Hello! I'm trying out smtp to send emails to recipients.")

		// err := smtp.SendMail("mail.teclu.com:25", a.smtpAuth, "jmiranda@teclu.com", mails, body.Bytes())
		// if err != nil{
		// 	log.Println(err)
		// }
	}()
}

// func (a *accountUseCase) RegisterCliente(ctx context.Context, form *account.RegisterForm) (res user.ClienteResponse, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
// 	defer cancel()

// 	res, err = a.accountRepo.RegisterCliente(ctx, form)
// 	return
// }

// func (a *accountUseCase) RegisterFuncionario(ctx context.Context, form *account.RegisterForm) (res user.UserAuth, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
// 	defer cancel()

// 	// res ,err =

// 	res, err = a.accountRepo.RegisterFuncionario(ctx, form)
// 	return
// }
