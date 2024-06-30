package database

import (
	"context"
)

type Database interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Put(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) error
	Post(ctx context.Context, key string, value []byte) error
}
