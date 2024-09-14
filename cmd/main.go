package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sarath/backend_project/cmd/api"
	"time"
)

func main() {
	var config api.Config

	// Reading the flags of the application
	flag.IntVar(&config.Port, "port", 4000, "API server port")
	flag.StringVar(&config.Env, "env", "dev", "Environment (dev | stag | production)")
  flag.StringVar(&config.Db.Dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgresSQL DSN")
  flag.StringVar(&config.JWTSecretKey, "jwt-secret", os.Getenv("jwt-secret"), "The JWT Secret key")
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
	//  db, err := OpenDB(app.Config)
	//  if err != nil{
	//    app.Logger.Fatal(err)
	//  }
	//  defer db.Close()
	//  app.Logger.Printf("database connection pool established")
	//
	// // start the server
	// app.Logger.Printf("Starting %s server on %s", app.Config.Env, server.Addr)
  err := server.ListenAndServe()
	app.Logger.Fatal(err)
}
