package api

import "log"

// Delcaring the version global constant
const version = "1.0.0"

type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn string // data source name
	}
}

type Application struct {
	Config Config
	Logger log.Logger
}
