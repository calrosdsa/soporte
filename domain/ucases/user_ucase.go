package ucases

import (
	"bytes"
	"context"
	"fmt"

	// "errors"
	// "fmt"
	"log"
	// "net/smtp"
	jwt "soporte-go/api/routes"
	"soporte-go/core/model/user"
	"time"

	"gopkg.in/gomail.v2"

	// "log"
	"html/template"
	"soporte-go/core/model"
)

type userUseCase struct {
	userRepo       user.UserRepository
	contextTimeout time.Duration
	gomailAuth     *gomail.Dialer
	util           model.Util
}

func NewUserUseCases(u user.UserRepository, timeout time.Duration, g *gomail.Dialer, util model.Util) user.UserUseCases {
	return &userUseCase{
		userRepo:       u,
		contextTimeout: timeout,
		gomailAuth:     g,
		util:           util,
	}
}

// func(a *userUseCase)GetUsersEmpresa(ctx context.Context,emId int,rol int)(res []user.UserForList,err error){
// 	ctx,cancel := context.WithTimeout(ctx,a.contextTimeout)
// 	defer cancel()
// 	if a.util.IsClienteAdmin(rol){
// 		res,err = a.userRepo.GetClientesEmpresa(ctx,emId)
// 	} else if a.util.IsFuncionarioAdmin(rol) {
// 		res,err = a.userRepo.GetFuncionariosEmpresa(ctx,emId)
// 	}
// 	return
// }

func (a *userUseCase) GetUsersbyEmpresaId(ctx context.Context,emId int) (res []user.UserForList,err error) {
	ctx, cancel := context.WithTimeout(ctx,a.contextTimeout)
	defer cancel()
	res,err = a.userRepo.GetClientesEmpresa(ctx,emId)
	return
}

func (a *userUseCase) GetUsersEmpresaByRol(ctx context.Context,emId int,rol int) (res []user.UserForList,err error) {
	ctx,cancel := context.WithTimeout(ctx,a.contextTimeout)
	defer cancel()
	if a.util.IsClienteAdmin(rol){
		res,err = a.userRepo.GetClientesEmpresa(ctx,emId)
	} else if a.util.IsFuncionarioAdmin(rol) {
		res,err = a.userRepo.GetFuncionariosEmpresa(ctx,emId)
	}
	return
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
		m.SetHeader("To", emails...)
		// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "Hello!")
		m.SetBody("text/html", body.String())
		// m.Attach("/home/Alex/lolcat.jpg")

		// d := gomail.NewDialer("mail.teclu.com", 25, "jmiranda@teclu.com", "jmiranda2022")
		if err := a.gomailAuth.DialAndSend(m); err != nil {
			log.Println("Error sending email",err)
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

func (a *userUseCase) ReSendEmail(m []string, url string) {
	a.sendEmail(m, url)
}

func (a *userUseCase) UserRegisterInvitation(ctx context.Context, to *user.UserRegistrationRequest,
	 id string, rol int, empresaId int) (res []user.UserShortInfo, err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	invitations := make([]user.UserShortInfo, len(to.To))
	for index, value := range to.To {
		tokenInvitation, _ := jwt.GenerateInvitationJWT(id, rol, empresaId, value)
		url := fmt.Sprintf("https://template-f-pearl.vercel.app/auth/registro?auth=%s", tokenInvitation)
		t := user.UserShortInfo{
			Nombre:  value,
			Id:      id,
			IsAdmin: to.IsAdmin,
		}
		if a.util.IsClienteRol(rol){
			val, err := a.userRepo.CreateUserInvitationC(ctx, &t)
			log.Println(err)
			invitations[index] = val
		}else if a.util.IsFuncionarioRol(rol){
			val, err := a.userRepo.CreateUserInvitationF(ctx, &t)
			log.Println(err)
			invitations[index] = val
		}
		// log.Println(err)
		// if err != nil{
		// 	return nil,err
		// }
		a.sendEmail(to.To, url)
	}

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

func (u *userUseCase) SearchUser(ctx context.Context, id string, q string) (res []user.UserShortInfo, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	res, err = u.userRepo.SearchUser(ctx, id, q)
	return
}

func (u *userUseCase) GetUsersShortIInfo(ctx context.Context, id string, rol int, emId int) (res []user.UserShortInfo, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	var list []user.UserShortInfo
	if u.util.IsFuncionarioAdmin(rol) {
		list, err = u.userRepo.GetUsersShortIInfoF(ctx,emId)
		if err != nil {
			return
		}
	} else if u.util.IsClienteAdmin(rol){
		list, err = u.userRepo.GetUsersShortIInfoC(ctx, id)
		if err != nil {
			return
		}
	}
	list2, err := u.userRepo.GetInvitaciones(ctx, id)
	res = append(list, list2...)
	return
}

func (u *userUseCase) GetClientesByArea(ctx context.Context, id int) (res []user.UserArea, err error) {
	res, err = u.userRepo.GetClientesByArea(ctx, id)
	return
}

func (u *userUseCase) GetUserAddList(ctx context.Context, f int, rol int, sId string) (res []user.UserArea, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	res, err = u.userRepo.GetUserAddList(ctx, f, rol, sId)
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

func (u *userUseCase) GetClientes(ctx context.Context, id string, rol int) ([]user.UserShortInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	list, err := u.userRepo.GetUsersShortIInfoC(ctx, id)
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
	err := u.userRepo.UpdateFuncionario(ctx, columns, values...)
	if err != nil {
		log.Println(err)
		return err
	}
	// query := `update clientes set nombre = &1,apellido = &2,email=&3,celular=&4,telefono=&5,profile_photo=&6)`
	return nil
}

func (u *userUseCase) GetUserById(ctx context.Context, id string, rol int) (res user.UserDetail, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	if u.util.IsClienteRol(rol){
		res, err = u.userRepo.GetClienteDetail(ctx, id)
		if err != nil {
			return
		}
	}else if u.util.IsFuncionarioRol(rol){
		res, err = u.userRepo.GetFuncionarioDetail(ctx, id)
		if err != nil {
			return
		}
	}
	return 
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
