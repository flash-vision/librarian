package librarian

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	Client *s3.S3
	Bucket string
}

func NewS3Storage() *S3Storage {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(viper.GetString("aws.region")),
	}))
	return &S3Storage{
		Client: s3.New(sess),
		Bucket: viper.GetString("aws.bucket"),
	}
}

func (s *S3Storage) Store(ctx context.Context, key string, value []byte) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(value
