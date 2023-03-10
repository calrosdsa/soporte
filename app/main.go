package main

import (
	// "bytes"
	"context"
	"fmt"
	"os"
	initApp "soporte-go/api"

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
	viper.SetConfigFile(`.env`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// if viper.GetBool(`debug`) {
	// 	log.Println("Service RUN on DEBUG mode")
	// }
}

// var auth smtp.Auth

// func SmtpEmial() {
// 	auth = smtp.PlainAuth("", "jorgemiranda0180@gmail.com", "opcpmdfaqrhtwwws", "smtp.gmail.com")
// 	t,_ := template.ParseFiles("templates/reset_password.html")
// 	var body bytes.Buffer
// 	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
// 	body.Write([]byte(fmt.Sprintf("Subject: yourSubject\n%s\n\n",headers)))
// 	t.Execute(&body,struct{
// 		Name  string
// 		URL  string
// 	}{
// 		Name: "AMAMAM",
// 		URL: "google.com",
// 	})
// 	smtp.SendMail("smtp.gmail.com:587",auth,"jorgemiranda0180@gmail.com",[]string{"alejandro12ab34cd@gmail.com"},body.Bytes())
// }

func main() {
	// SmtpEmial()
	
	creds := credentials.NewStaticCredentials(viper.GetString("AWS_ID"), viper.GetString("AWS_SECRET"), "")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("sa-east-1"),
		Credentials: creds,
	})
	if err != nil {
		exitErrorf("%v", err)
	}
	db, err := pgxpool.New(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	// defer db.Close()

	err = db.Ping(context.Background())
	CheckError(err)
	fmt.Println("Connected!")
	initApp.InitServer(db, context.Background(),sess)

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
