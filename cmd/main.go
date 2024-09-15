package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"

	"sarath/backend_project/cmd/api"
	"sarath/backend_project/internal/cache"
	"sarath/backend_project/internal/data"
	filestore "sarath/backend_project/internal/file_store"
)

func main() {
	var config api.Config

	// Reading the flags of the application
	flag.IntVar(&config.Port, "port", 4000, "API server port")
	flag.StringVar(&config.Env, "env", "dev", "Environment (dev | stag | production)")
	flag.StringVar(&config.Db.Dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgresSQL DSN")
	flag.StringVar(&config.Jwt.Secret, "jwt-secret", os.Getenv("JWT_SECRET"), "The JWT Secret key")
	flag.StringVar(&config.Aws.Bucket, "aws-bucket", os.Getenv("AWS_BUCKET"), "The AWS Bucket which the files are uploaded")
	flag.Parse()

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := aws.String(os.Getenv("AWS_REGION"))
	awsToken := ""

	redisAddr := os.Getenv("REDIS_ADDRESS")
	// redisPassword := os.Getenv("REDIS_PASSWORD")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// starting the db
	db, err := OpenDB(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection pool established")

	// creating a aws session
	awsSess, err := session.NewSession(&aws.Config{
		Region:      awsRegion,
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, awsToken),
	})
	if err != nil {
		logger.Fatal(err)
	}

	// creating a redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		// Password: redisPassword,
		DB: 0,
	})

	// ping the redis client
	_, err = redisClient.Ping().Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("redis connection pool established")

	// defining the application
	app := &api.Application{
		Config: config,
		Logger: logger,
		Models: data.NewModels(db),
		FileStore: &filestore.AWSFileStore{
			Bucket: config.Aws.Bucket,
			S3Sess: awsSess,
		},
		Cache: &cache.RedisCache{
			Client: redisClient,
		},
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	// run the background delete from db task
	app.Background(func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for v := range ticker.C {
			app.Logger.Printf("Running delete from db task at %v", v)
			app.DeleteFromDBTask()
		}
	})

	// start the server
	app.Logger.Printf("Starting %s server on %s", app.Config.Env, server.Addr)
	err = server.ListenAndServe()
	app.Logger.Fatal(err)
}

func OpenDB(cfg api.Config) (*sql.DB, error) {
	// open an sql connection
	db, err := sql.Open("postgres", cfg.Db.Dsn)
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
