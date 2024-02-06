package librarian

import "context"

type Storer interface {
	Store(ctx context.Context, key string, value []byte) error
	Retrieve(ctx context.Context, key string) ([]byte, error)
}
