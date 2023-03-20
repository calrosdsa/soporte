package api

import (
	"context"
	"log"
	"time"

	// "fmt"
	// "net/smtp"

	_ "github.com/lib/pq"
	"gopkg.in/gomail.v2"

	// "github.com/spf13/viper"

	// "github.com/spf13/viper"

	// "github.com/spf13/viper"
	_r_ws "soporte-go/api/routes"
	_r_account "soporte-go/api/routes/account"
	_r_caso "soporte-go/api/routes/caso"
	_r_empresa "soporte-go/api/routes/empresa"
	_r_media "soporte-go/api/routes/media"
	_r_user "soporte-go/api/routes/user"

	_repository "soporte-go/core/repository"
	_uCase "soporte-go/domain/ucases"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitServer(db *pgxpool.Pool, ctx context.Context, sess *session.Session) {
	e := echo.New()
	log.Println("init server....`")
	// middl := InitMiddleware()
	// e.Use(middl.CORS)
	// e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		// AllowMethods: []string{"*"},
	}))
	// e.Use(middleware.Recover())
	mUser := "jmiranda@teclu.com"
	mPass := "jmiranda2022"
	mHost := "mail.teclu.com"
	gomailAuth := gomail.NewDialer(mHost, 25, mUser, mPass)

	timeoutContext := time.Duration(15) * time.Second
	accountRepository := _repository.NewPgAccountRepository(db, ctx)
	accountUseCase := _uCase.NewAccountUseCase(accountRepository, timeoutContext)
	_r_account.NewAccountHandler(e, accountUseCase)

	//user
	userRepository := _repository.NewPgUserRepository(db, ctx)
	userUseCase := _uCase.NewUserUseCases(userRepository, timeoutContext, gomailAuth)
	_r_user.NewUserHandler(e, userUseCase)

	//caso
	casoRepository := _repository.NewPgCasoRepository(db, ctx)
	casoUseCase := _uCase.NewCasoUseCase(casoRepository, timeoutContext)
	_r_caso.NewCasoHandler(e, casoUseCase)

	//empresa
	empresaRepository := _repository.NewPgEmpresaRepository(db, ctx)
	empresaUseCase := _uCase.NewEmpresaUseCase(empresaRepository, timeoutContext)
	_r_empresa.NewEmpresaHandler(e, empresaUseCase)

	//media
	mediaRepository := _repository.NewMediaRepository(db, ctx)
	mediaUseCase := _uCase.NewMediaUseCase(mediaRepository, timeoutContext, sess)
	_r_media.NewMediaHandler(e, mediaUseCase)

	wsRepository := _repository.NewWsRepository(db, ctx)
	_r_ws.NewWsHandler(e, wsRepository)
	log.Fatal(e.Start(":8000")) //nolint

}
