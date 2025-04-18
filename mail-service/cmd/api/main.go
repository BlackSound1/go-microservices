package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const WEB_PORT = "80"

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail service on port", WEB_PORT)

	srv := &http.Server{
		Addr:    ":" + WEB_PORT,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// createMail reads in environment variables and creates a Mail object based on them
func createMail() Mail {

	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
	}

	return m
}
