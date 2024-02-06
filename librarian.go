package librarian

import (
	"context"
	"log"
	"os"
	"sync"
)

// Librarian orchestrates storage and retrieval across multiple backends.
type Librarian struct {
	Memcached  Storer
	InMemory   *InMemoryStorage
	logger     *log.Logger
}

func NewLibrarian(memcached Storer, logFilePath string) (*Librarian, error) {
	logger, err := NewLogger(logFilePath)
	if err != nil {
		return nil, err
	}

	return &Librarian{
		Memcached: memcached,
		InMemory:  NewInMemoryStorage(),
		logger:    logger,
	}, nil
}

func (l *Librarian) Store(ctx context.Context, key string, value []byte) error {
	// Store in Memcached
	if err := l.Memcached.Store(ctx, key, value); err != nil {
		l.logger.Printf("Error storing key %s in Memcached: %v", key, err)
		return err
	}

	// Store in InMemory
	l.InMemory.Store(key, value)
	l.logger.Printf("Stored key %s", key)

	return nil
}

func (l *Librarian) Retrieve(ctx context.Context, key string) ([]byte, error) {
	// Attempt to retrieve from InMemory first
	if value, found := l.InMemory.Retrieve(key); found {
		l.logger.Printf("Retrieved key %s from InMemory", key)
		return value, nil
	}

	// Fallback to Memcached
	value, err := l.Memcached.Retrieve(ctx, key)
	if err != nil {
		l.logger.Printf("Error retrieving key %s from Memcached: %v", key, err)
		return nil, err
	}
	l.logger.Printf("Retrieved key %s from Memcached", key)
	return value, nil
}
