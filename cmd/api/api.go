package api

import (
	"log"

	"sarath/backend_project/internal/cache"
	"sarath/backend_project/internal/data"
	filestore "sarath/backend_project/internal/file_store"
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
	Config    Config
	Logger    *log.Logger
	Models    *data.Models
	FileStore *filestore.FileStore
	Cache     *cache.Cache
}

func (app *Application) Background(fn func()) {
	// Launch a background goroutine.
	go func() {
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
        app.Logger.Printf("Recovered from a panic: %v", err)
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}
