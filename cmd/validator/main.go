package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ArminGh02/go-auth-system/internal/config"
	"github.com/ArminGh02/go-auth-system/internal/imagga"
	"github.com/ArminGh02/go-auth-system/internal/mailgun"
	"github.com/ArminGh02/go-auth-system/internal/mongodb"
	"github.com/ArminGh02/go-auth-system/internal/rabbitmq"
	"github.com/ArminGh02/go-auth-system/internal/s3"
	"github.com/ArminGh02/go-auth-system/internal/svc/validator"
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
	log.Println("current configurations:", cfg)

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

	imagga := imagga.New(cfg.Imagga)

	mg := mailgun.New(cfg.MailGun)

	v := validator.NewListener(mongo, s3, q, imagga, mg)

	log.Println("validator started listening to registrations queue")
	return v.Listen()
}
