package api

import (
	"context"
	"database/sql"
	"log"
	"time"

	// "fmt"
	// "net/smtp"

	_ "github.com/lib/pq"
	"gopkg.in/gomail.v2"

	// "github.com/spf13/viper"

	// "github.com/spf13/viper"

	// "github.com/spf13/viper"
	_account "soporte-go/api/routes/account"
	_caso "soporte-go/api/routes/caso"
	_empresa "soporte-go/api/routes/empresa"
	_media "soporte-go/api/routes/media"
	_user "soporte-go/api/routes/user"
	_ws "soporte-go/api/routes/ws"

	domain_r "soporte-go/domain/repository"

	_r_caso "soporte-go/core/repository/caso"

	// "soporte-go/core/model"
	_repository "soporte-go/core/repository"
	_uCase "soporte-go/domain/ucases"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitServer(db *pgxpool.Pool, db2 *sql.DB, ctx context.Context, sess *session.Session) {
	e := echo.New()
	log.Println("init server....`")
	// middl := InitMiddleware()
	// e.Use(middl.CORS)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	e.Use(middleware.Recover())
	// e.Use(middleware.Logger())
	mUser := "jmiranda@teclu.com"
	mPass := "jmiranda2022"
	mHost := "mail.teclu.com"
	gomailAuth := gomail.NewDialer(mHost, 25, mUser, mPass)

	util := domain_r.NewUtil()

	timeoutContext := time.Duration(15) * time.Second
	accountRepository := _repository.NewPgAccountRepository(db2, ctx)
	accountUseCase := _uCase.NewAccountUseCase(accountRepository, timeoutContext, util,gomailAuth)
	_account.NewAccountHandler(e, accountUseCase)

	utilRepository := _repository.NewPgUtilRepository(db2)
	//user
	userRepository := _repository.NewPgUserRepository(db2, ctx)
	userUseCase := _uCase.NewUserUseCases(userRepository, timeoutContext, gomailAuth, util)
	_user.NewUserHandler(e, userUseCase)

	//caso
	casoRepository := _r_caso.NewPgCasoRepository(db2, ctx)
	casoUseCase := _uCase.NewCasoUseCase(casoRepository, timeoutContext, util,gomailAuth,utilRepository)
	_caso.NewCasoHandler(e, casoUseCase)

	
	
	//empresa
	empresaRepository := _repository.NewPgEmpresaRepository(db2, ctx)
	empresaUseCase := _uCase.NewEmpresaUseCase(empresaRepository, timeoutContext, util)
	_empresa.NewEmpresaHandler(e, empresaUseCase)

	//caso
	casoDataRepository := _r_caso.NewCasoDataRepository(db2)
	casoDataUseCase := _uCase.NewCasoDataUseCase(casoDataRepository,empresaRepository, timeoutContext)
	_caso.NewCasoDataHandler(e, casoDataUseCase)

	//media
	mediaRepository := _repository.NewMediaRepository(db2, ctx)
	mediaUseCase := _uCase.NewMediaUseCase(mediaRepository, timeoutContext, sess)
	_media.NewMediaHandler(e, mediaUseCase)

	wsRepository := _repository.NewWsRepository(db2, ctx)
	_ws.NewWsHandler(e, wsRepository)
	log.Fatal(e.Start(":8000")) //nolint

}
