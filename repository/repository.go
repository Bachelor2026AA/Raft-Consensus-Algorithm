package repository

import (
	"context"

	bolt "go.etcd.io/bbolt"
)

type LogEntryRepository interface {
	append(ctx context.Context, index int) error
	get(ctx context.Context) error
	getFromIndex(ctx context.Context) error
	deleteFromIndex(ctx context.Context) error
}

type LogEntryRepositoryImpl struct {
	db *bolt.DB
}
