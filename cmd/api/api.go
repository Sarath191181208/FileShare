package api

import (
	"log"

	"sarath/backend_project/internal/data"

	"github.com/aws/aws-sdk-go/aws/session"
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
	Aws struct {
		Bucket string
	}
}

type Application struct {
	Config Config
	Logger *log.Logger
	Models *data.Models
	S3Sess *session.Session
}
