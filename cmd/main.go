package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"sarath/backend_project/cmd/api"
)

func main() {
	var config api.Config

	// Reading the flags of the application
	flag.IntVar(&config.Port, "port", 4000, "API server port")
	flag.StringVar(&config.Env, "env", "dev", "Environment (dev | stag | production)")
	flag.StringVar(&config.Db.Dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgresSQL DSN")
	flag.StringVar(&config.Jwt.Secret, "jwt-secret", os.Getenv("jwt-secret"), "The JWT Secret key")
	flag.Parse()

	// defining the application
	app := &api.Application{
		Config: config,
		Logger: *log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	// starting the db
	db, err := OpenDB(app.Config)
	if err != nil {
		app.Logger.Fatal(err)
	}
	defer db.Close()
	app.Logger.Printf("database connection pool established")

	// start the server
	app.Logger.Printf("Starting %s server on %s", app.Config.Env, server.Addr)
	err := server.ListenAndServe()
	app.Logger.Fatal(err)
}

func OpenDB(cfg Config) (*sql.DB, error) {
	// open an sql connection
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// create a session for the db
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check for invalid ping
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
