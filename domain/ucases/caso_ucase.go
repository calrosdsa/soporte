package ucases

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"soporte-go/core/model"
	"soporte-go/core/model/caso"
	"soporte-go/core/model/user"
	_util "soporte-go/core/model/util"
	"soporte-go/core/reportes/excel"
	"soporte-go/core/reportes/html"
	"soporte-go/core/reportes/pdf"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"gopkg.in/gomail.v2"
)

type casoUseCase struct {
	casoRepo       caso.CasoRepository
	contextTimeout time.Duration
	util           model.Util
	gomailAuth     *gomail.Dialer
	utilRepo _util.UtilRepository
}

func NewCasoUseCase(uc caso.CasoRepository, timeout time.Duration, util model.Util,
	m *gomail.Dialer,utilRepo _util.UtilRepository) caso.CasoUseCase {
	return &casoUseCase{
		casoRepo:       uc,
		contextTimeout: timeout,
		util:           util,
		gomailAuth: m,
		utilRepo: utilRepo,
	}
}


func (uc *casoUseCase) GetUsuariosCaso(ctx context.Context, cId string) (res []user.UserForList, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	res, err = uc.casoRepo.GetUsuariosCaso(ctx, cId)
	return
}

func (uc *casoUseCase) AsignarFuncionarioSoporte(ctx context.Context, u *caso.UserCaso) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	for _, id := range u.UserId {
		err = uc.casoRepo.AsignarFuncionarioSoporte(ctx, u.CasoId, id)
	}
	return
}

func (uc *casoUseCase) GetReporteCaso(ctx context.Context, t model.FileType, c caso.Caso) (b bytes.Buffer, err error) {
	var buffer bytes.Buffer
	res, err := uc.casoRepo.GetUsuariosCaso(ctx, c.Id)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	ms, err := uc.casoRepo.GetMessagesCaso(ctx, c.Id)
	if err != nil {
		log.Println(err)
	}
	log.Println(ms)
	switch t {
	case model.HTML:
		html.HtmlCasoReporte(&buffer, c, res, ms)
	}

	return buffer, err
}

func (uc *casoUseCase) GetReporteCasos(ctx context.Context, t model.FileType, options *caso.CasoReporteOptions) (b bytes.Buffer, err error) {
	var buffer bytes.Buffer
	casos, err := uc.casoRepo.GetCasosCliForReporte(ctx, options)
	if err != nil {
		log.Println(err)
	}
	casos2, err := uc.casoRepo.GetCasosFunForReporte(ctx, options)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(casos)
	switch t {
	case model.XLSX:
		err = excel.ReporteCasosExcel(casos, casos2, &buffer)
		if err != nil {
			log.Println(err)
			return
		}
	case model.PDF:
		err = pdf.ReporteCasos(casos, &buffer)
		if err != nil {
			log.Println(err)
			return
		}
	}

	return buffer, err
}

func (uc *casoUseCase) FinalizarCaso(ctx context.Context, fD *caso.FinalizacionDetail) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.casoRepo.FinalizarCaso(ctx, fD)
	return
}

func (uc *casoUseCase) AsignarFuncionario(ctx context.Context, id string, idF string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err,mail := uc.casoRepo.AsignarFuncionario(ctx, id, idF)
	log.Println(mail)
	// uc.sendEmail([]string{mail},"")
	return
}

func (uc *casoUseCase) GetCaso(ctx context.Context, id string, rol int) (res caso.Caso, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	if uc.util.IsClienteRol(rol) {
		res, err = uc.casoRepo.GetCasoCliente(ctx, id)
		if err != nil {
			return
		}
	} else if uc.util.IsFuncionarioRol(rol) {
		res, err = uc.casoRepo.GetCasoFuncionario(ctx, id)
		if err != nil {
			return
		}
	}
	return
}

func (uc *casoUseCase) GetCasosUser(ctx context.Context, id string, query *caso.CasoQuery, rol int) (res []caso.Caso, size int, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	page, offset := uc.util.PaginationValues(query.Page, query.PageSize)
	query.Page = page
	query.PageSize = offset
	// log.Println(query.Page,query.PageSize)
	if uc.util.IsClienteRol(rol) {
		res, err = uc.casoRepo.GetCasosCliente(ctx, id, query)
		if err != nil {
			return
		}
		size, err = uc.casoRepo.GetCasosCountCliente(ctx, id)
	} else if uc.util.IsFuncionarioRol(rol) {
		res, err = uc.casoRepo.GetCasosFuncionario(ctx, id, query)
		if err != nil {
			return
		}
		size, err = uc.casoRepo.GetCasosCountFuncionario(ctx, id)
	}
	return
}

func (uc *casoUseCase) GetAllCasosUser(ctx context.Context, id string, query *caso.CasoQuery, rol int) (res []caso.Caso, size int, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	page, offset := uc.util.PaginationValues(query.Page, 30)
	query.Page = page
	query.PageSize = offset
	if uc.util.IsClienteAdmin(rol) {
		size, err = uc.casoRepo.GetCasosCountbySuperiorId(ctx, id)
		if err != nil {
			return
		}
		res, err = uc.casoRepo.GetAllCasosUserCliente(ctx, id, query)
	} else if uc.util.IsFuncionarioAdmin(rol) {
		size, err = uc.casoRepo.GetCasosCount(ctx)
		if err != nil {
			return
		}
		res, err = uc.casoRepo.GetAllCasosUserFuncionario(ctx, 0, query)
	}
	return
}

func (uc *casoUseCase) GetCasosFromUserCaso(ctx context.Context, id string, q *caso.CasoQuery) (res []caso.Caso, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	page, offset := uc.util.PaginationValues(q.Page, 30)
	q.Page = page
	q.PageSize = offset
	res, err = uc.casoRepo.GetCasosFromUserCaso(ctx, id, q)
	return
}

func (uc *casoUseCase) UpdateCaso(ctx context.Context, c *caso.Caso) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	err = uc.casoRepo.UpdateCaso(ctx, c)
	return
}

func (uc *casoUseCase) CreateCaso(ctx context.Context, cas *caso.Caso, id string, emI int, rol int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	log.Println(rol)
	if uc.util.IsClienteRol(rol) {
		err = uc.casoRepo.CreateCasoCliente(ctx, cas, id, emI, rol)
		if err != nil {
			log.Println(err)
			return
		}
		go func() {
			ctxB := context.Background()
			ctx,cancel  := context.WithTimeout(ctxB,uc.contextTimeout)
			defer cancel()
			d,_ := uc.utilRepo.GetEmailFromUserCasos(ctx,*cas.Area,id)
			m:= make([]string,len(d.Mails))
			for i,v := range d.Mails {
				m[i] = v.Mail
			}
			uc.sendEmail(m,fmt.Sprintf(`%s/casos/%s?r=%s`, viper.GetString("CLIENT_URL"),cas.Id,strconv.Itoa(rol)),
			d.UsuarioName + " "+ d.UsuarioApellido,d.Area,d.Entidad)
		}()
	} else if uc.util.IsFuncionarioRol(rol) {
		err = uc.casoRepo.CreateCasoFuncionario(ctx, cas, id, emI, rol)
		if err != nil {
			return
		}
	}
	return
}


func (uc *casoUseCase) sendEmail(emails []string, url string,user string,area string,entidad string) {
	log.Println(emails)
	// mails := []string{"diegoarmando12ab34cd@gmail.com"}
	go func() {
		t, _ := template.ParseFiles("templates/new_caso.html")
		var body bytes.Buffer
		t.Execute(&body, struct {
			URL string
			USUARIO string
			AREA string
			ENTIDAD string
		}{
			URL: url,
			USUARIO:user,
			AREA:area,
			ENTIDAD:entidad,
		})
		m := gomail.NewMessage()
	
		m.SetHeader("From", "jmiranda@teclu.com")
		m.SetHeader("To", emails...)
		// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "Se ha creado un nuevo caso")
		m.SetBody("text/html", body.String())
		// m.Attach("/home/Alex/lolcat.jpg")

		// d := gomail.NewDialer("mail.teclu.com", 25, "jmiranda@teclu.com", "jmiranda2022")
		if err := uc.gomailAuth.DialAndSend(m); err != nil {
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