package librarian

import (
	"context"
	"errors"
	"sync"
)

type InMemoryStorage struct {
	store map[string][]byte
	mu    sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		store: make(map[string][]byte),
	}
}

func (ims *InMemoryStorage) Store(ctx context.Context, key string, value []byte) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.store[key] = value
	return nil
}

func (ims *InMemoryStorage) Retrieve(ctx context.Context, key string) ([]byte, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	value, found := ims.store[key]
	if !found {
		return nil, errors.New("key not found")
	}
	return value, nil
}
