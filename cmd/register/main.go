package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ArminGh02/go-auth-system/internal/config"
	"github.com/ArminGh02/go-auth-system/internal/mongodb"
	"github.com/ArminGh02/go-auth-system/internal/rabbitmq"
	"github.com/ArminGh02/go-auth-system/internal/s3"
	"github.com/ArminGh02/go-auth-system/internal/svc/register"
)

const timeout = 10 * time.Second

func main() {
	log.Fatal(run())
}

func run() (err error) {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	log.Printf("current configurations: %+v\n", cfg)

	ctx, done := context.WithTimeout(context.Background(), timeout)
	defer done()
	mongo, err := mongodb.New(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer func() {
		ctx, done := context.WithTimeout(context.Background(), timeout)
		defer done()
		err = errors.Join(err, mongo.Close(ctx))
	}()
	log.Println("connected to mongodb successfully")

	q, err := rabbitmq.New(cfg.Broker)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, q.Close())
	}()
	log.Println("connected to rabbitmq successfully")

	s3, err := s3.New(cfg.S3)
	if err != nil {
		return err
	}
	log.Println("connected to s3 successfully")

	h := register.NewHandler(mongo, s3, q)

	log.Println("register service is running on port:", cfg.ServerPort)
	return http.ListenAndServe(":"+cfg.ServerPort, h)
}
