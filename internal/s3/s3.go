package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	Config  Config
	session *session.Session
	client  *s3.S3
}

func New(cfg Config) (*S3, error) {
	session, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(cfg.Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		Region:      aws.String(cfg.Region),
	})
	if err != nil {
		return nil, err
	}

	return &S3{
		Config:  cfg,
		session: session,
		client:  s3.New(session),
	}, nil
}

func (s *S3) Put(ctx context.Context, key string, file io.ReadSeeker) error {
	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(key),
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
	})
	return err
}

func (s *S3) URL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.Config.Endpoint, s.Config.Bucket, key)
}
