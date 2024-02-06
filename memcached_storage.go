package librarian

import (
	"context"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

// MemcachedStorage implements the Storer interface for Memcached.
type MemcachedStorage struct {
	client *memcache.Client
}

// NewMemcachedStorage creates a new instance of MemcachedStorage.
func NewMemcachedStorage(serverAddresses []string) *MemcachedStorage {
	client := memcache.New(serverAddresses...)
	return &MemcachedStorage{client: client}
}

// Store saves a key-value pair in Memcached.
func (m *MemcachedStorage) Store(ctx context.Context, key string, value []byte) error {
	// Create a Memcached item
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(time.Hour.Seconds()), // Example expiration, adjust as needed
	}

	// Use context deadline if available
	if deadline, ok := ctx.Deadline(); ok {
		m.client.Timeout = time.Until(deadline)
	}

	// Set the item in Memcached
	err := m.client.Set(item)
	if err != nil {
		return fmt.Errorf("memcached store error: %w", err)
	}
	return nil
}

// Retrieve fetches a value by key from Memcached.
func (m *MemcachedStorage) Retrieve(ctx context.Context, key string) ([]byte, error) {
	// Use context deadline if available
	if deadline, ok := ctx.Deadline(); ok {
		m.client.Timeout = time.Until(deadline)
	}

	item, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, fmt.Errorf("memcached retrieve error: key not found")
		}
		return nil, fmt.Errorf("memcached retrieve error: %w", err)
	}
	return item.Value, nil
}
