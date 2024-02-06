package librarian

import (
	"context"
	"cloud.google.com/go/storage"
	"io/ioutil"
	"github.com/spf13/viper"
)

// GCSStorage implements the Storer interface for Google Cloud Storage.
type GCSStorage struct {
	client    *storage.Client
	bucketName string
	keyPrefix  string // Added keyPrefix to specify the beginning of every key
}

// NewGCSStorage creates a new instance of GCSStorage.
func NewGCSStorage(ctx context.Context) (*GCSStorage, error) {
	// Initialize Viper to read the configuration
	viper.SetConfigName("config") // Configuration file name without the extension
	viper.SetConfigType("yml")    // or viper.SetConfigType("YAML")
	viper.AddConfigPath(".")      // Path to look for the configuration file in

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Extract GCS configuration
	bucketName := viper.GetString("gcs.bucketName")
	keyPrefix := viper.GetString("gcs.keyPrefix")

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCSStorage{
		client:     client,
		bucketName: bucketName,
		keyPrefix:  keyPrefix,
	}, nil
}

func (g *GCSStorage) Store(ctx context.Context, key string, value []byte) error {
	bucket := g.client.Bucket(g.bucketName)
	object := bucket.Object(g.keyPrefix + key) // Prepend the keyPrefix to the key
	writer := object.NewWriter(ctx)

	if _, err := writer.Write(value); err != nil {
		writer.Close()
		return err
	}

	return writer.Close()
}

func (g *GCSStorage) Retrieve(ctx context.Context, key string) ([]byte, error) {
	bucket := g.client.Bucket(g.bucketName)
	object := bucket.Object(g.keyPrefix + key) // Prepend the keyPrefix to the key
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
