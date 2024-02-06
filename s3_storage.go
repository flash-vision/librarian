package librarian

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

// S3Storage implements the Storer interface for AWS S3.
type S3Storage struct {
	client     *s3.Client
	bucketName string
	keyPrefix  string
}

// NewS3Storage creates a new instance of S3Storage.
func NewS3Storage(ctx context.Context) (*S3Storage, error) {
	// Load configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	bucketName := viper.GetString("s3.bucketName")
	keyPrefix := viper.GetString("s3.keyPrefix")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		client:     s3.NewFromConfig(cfg),
		bucketName: bucketName,
		keyPrefix:  keyPrefix,
	}, nil
}

func (s *S3Storage) Store(ctx context.Context, key string, value []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s.keyPrefix + key),
		Body:   bytes.NewReader(value),
	})
	return err
}

func (s *S3Storage) Retrieve(ctx context.Context, key string) ([]byte, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s.keyPrefix + key),
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
