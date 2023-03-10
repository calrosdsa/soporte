package api

import (
	"context"
	"log"
	"net/smtp"
	"time"

	// "fmt"
	// "net/smtp"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	// "github.com/spf13/viper"

	// "github.com/spf13/viper"
	_router "soporte-go/api/routes"
	_repository "soporte-go/core/repository"
	_uCase "soporte-go/domain/ucases"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)





func InitServer(db *pgxpool.Pool,ctx context.Context,sess *session.Session){
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
		mUser := viper.GetString("EMAIL_USER")
		mPass := viper.GetString("EMAIL_PASSWORD")
		mHost := viper.GetString("EMAIL_HOST")
	    sAuth := smtp.PlainAuth("",mUser, mPass,mHost)

		timeoutContext := time.Duration(15) * time.Second
		accountRepository := _repository.NewPgAccountRepository(db,ctx)
		accountUseCase := _uCase.NewAccountUseCase(accountRepository,timeoutContext)
		_router.NewAccountHandler(e,accountUseCase)		

		//user
		userRepository := _repository.NewPgUserRepository(db,ctx)
		userUseCase := _uCase.NewUserUseCases(userRepository,timeoutContext,sAuth)
		_router.NewUserHandler(e,userUseCase)

		//caso
		casoRepository := _repository.NewPgCasoRepository(db,ctx)
		casoUseCase := _uCase.NewCasoUseCase(casoRepository,timeoutContext)
		_router.NewCasoHandler(e,casoUseCase)

		//empresa
		empresaRepository := _repository.NewPgEmpresaRepository(db,ctx)
		empresaUseCase := _uCase.NewEmpresaUseCase(empresaRepository,timeoutContext)
		_router.NewEmoresaHandler(e,empresaUseCase)

		//media
		mediaRepository := _repository.NewMediaRepository(db,ctx)
		mediaUseCase :=  _uCase.NewMediaUseCase(mediaRepository,timeoutContext,sess)
		_router.NewMediaHandler(e,mediaUseCase)

		wsRepository := _repository.NewWsRepository(db,ctx)
		_router.NewWsHandler(e,wsRepository)
    	log.Fatal(e.Start(":8000")) //nolint


}


