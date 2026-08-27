package repository

import (
	"context"

	bolt "go.etcd.io/bbolt"
)

type LogEntryRepository interface {
	Append(ctx context.Context, command string) error
	Get(ctx context.Context) error
	GetFromIndex(ctx context.Context, index int) error
	DeleteFromIndex(ctx context.Context) error
}

type LogEntryRepositoryImpl struct {
	db *bolt.DB
}
