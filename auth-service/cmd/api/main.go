package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BlackSound1/go-microservices/auth/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const WEB_PORT = "80"

var counts int64

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

func main() {
	log.Println("Starting auth service on port ", WEB_PORT)

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres")
	}

	app := Config{
		Client: &http.Client{},
	}

	srv := &http.Server{
		Addr:    ":" + WEB_PORT,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// openDB attempts to open a database pool with the given DSN string.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Test DB before returning it
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB attempts to establish a connection to the Postgres database using the DSN environment variable.
// It continuously tries (every 2 seconds) to connect until successful or the retry limit is exceeded.
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Trying again in 2 seconds...")

		time.Sleep(2 * time.Second)

		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
