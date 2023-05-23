package main

import (
	// "bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	initApp "soporte-go/api"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	// "net/smtp"
	// "soporte-go/util"
)

func init() {
	// viper.SetConfigFile(`/home/rootuser/soporte/app/.env`)
	viper.SetConfigFile(`.env`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// if viper.GetBool(`debug`) {
	// 	log.Println("Service RUN on DEBUG mode")
	// }
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12ab34cd56ef"
	dbname   = "soporte"
)

func main() {
	loc, err := time.LoadLocation("America/La_Paz")
	if err != nil {
		log.Println(loc)
	}
    time.Local = loc
	// SmtpEmial()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	creds := credentials.NewStaticCredentials(viper.GetString("AWS_ID"), viper.GetString("AWS_SECRET"), "")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("sa-east-1"),
		Credentials: creds,
	})
	if err != nil {
		exitErrorf("%v", err)
	}
	db2, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	db, err := pgxpool.New(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}


	// defer db.Close()

	CheckError(err)
	fmt.Println("Connected!")
	initApp.InitServer(db,db2, context.Background(),sess)

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
