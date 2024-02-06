package librarian

import (
	"context"
	"errors"
	"sync"
)

// InMemoryStorage provides a basic in-memory key-value store.
type InMemoryStorage struct {
	store map[string][]byte
	mu    sync.RWMutex
}

// NewInMemoryStorage initializes a new InMemoryStorage instance.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		store: make(map[string][]byte),
	}
}

// Store saves a key-value pair in the memory.
// Adjusted to match the Storer interface if it requires context.Context
func (ims *InMemoryStorage) Store(ctx context.Context, key string, value []byte) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.store[key] = value
	return nil // Assuming storing in memory always succeeds
}

// Retrieve fetches a value by key from the memory.
// Modified to meet the Storer interface requirements
func (ims *InMemoryStorage) Retrieve(ctx context.Context, key string) ([]byte, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	value, found := ims.store[key]
	if !found {
		return nil, errors.New("key not found")
	}
	return value, nil
}
