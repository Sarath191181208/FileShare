package api

import (
	"log"
	"sarath/backend_project/internal/data"
)

// Delcaring the version global constant
const version = "1.0.0"

type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn string // data source name
	}
	Jwt struct {
		Secret string
	}
}

type Application struct {
	Config Config
	Logger *log.Logger
  Models *data.Models
}
