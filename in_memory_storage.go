package librarian

import "sync"

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
func (ims *InMemoryStorage) Store(key string, value []byte) {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.store[key] = value
}

// Retrieve fetches a value by key from the memory.
func (ims *InMemoryStorage) Retrieve(key string) ([]byte, bool) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	value, found := ims.store[key]
	return value, found
}
