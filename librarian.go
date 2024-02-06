package librarian

import (
	"context"
	"errors"
	"log"
	"os"
)

// StorageType defines different types of storage backends.
type StorageType int

const (
	CloudStorageType StorageType = iota
	PostgresType
	MemcachedType
	InMemoryType
)

// StorageBackend associates a storage type with its implementation.
type StorageBackend struct {
	Type   StorageType
	Storer Storer
}

// Librarian orchestrates storage and retrieval across multiple backends.
type Librarian struct {
	backends map[StorageType]Storer
	logger   *log.Logger
}

// NewLibrarian initializes a new Librarian instance with specified backends.
// It returns an error if no valid storage backend is provided or if there's an issue with the logger setup.
func NewLibrarian(backends []StorageBackend, logFilePath string) (*Librarian, error) {
	if len(backends) == 0 {
		return nil, errors.New("no storage backends provided")
	}

	// Setup logger
	loggerFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	logger := log.New(loggerFile, "LIBRARIAN: ", log.Ldate|log.Ltime|log.Lshortfile)

	librarian := &Librarian{
		backends: make(map[StorageType]Storer),
		logger:   logger,
	}

	// Initialize backends
	for _, backend := range backends {
		if backend.Storer == nil {
			return nil, errors.New("a provided storage backend is not initialized")
		}
		librarian.backends[backend.Type] = backend.Storer
	}

	return librarian, nil
}

func (l *Librarian) Store(ctx context.Context, key string, value []byte) error {
	var lastErr error
	for _, backend := range l.backends {
		if err := backend.Store(ctx, key, value); err != nil {
			l.logger.Printf("Error storing key %s in backend: %v", key, err)
			lastErr = err
		}
	}
	return lastErr // Returns the last error encountered, if any
}

func (l *Librarian) Retrieve(ctx context.Context, key string) ([]byte, error) {
	for _, backend := range l.backends {
		if value, err := backend.Retrieve(ctx, key); err == nil {
			return value, nil
		} else {
			l.logger.Printf("Error retrieving key %s from backend: %v", key, err)
		}
	}
	return nil, errors.New("key not found in any backend")
}
