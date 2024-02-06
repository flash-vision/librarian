package librarian

import (
	"context"
	"cloud.google.com/go/storage"
	"io/ioutil"
)

// GCSStorage implements the Storer interface for Google Cloud Storage.
type GCSStorage struct {
	client *storage.Client
	bucketName string
}

// NewGCSStorage creates a new instance of GCSStorage.
func NewGCSStorage(ctx context.Context, bucketName string) (*GCSStorage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCSStorage{
		client: client,
		bucketName: bucketName,
	}, nil
}

func (g *GCSStorage) Store(ctx context.Context, key string, value []byte) error {
	bucket := g.client.Bucket(g.bucketName)
	object := bucket.Object(key)
	writer := object.NewWriter(ctx)

	if _, err := writer.Write(value); err != nil {
		writer.Close()
		return err
	}

	return writer.Close()
}

func (g *GCSStorage) Retrieve(ctx context.Context, key string) ([]byte, error) {
	bucket := g.client.Bucket(g.bucketName)
	object := bucket.Object(key)
	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
