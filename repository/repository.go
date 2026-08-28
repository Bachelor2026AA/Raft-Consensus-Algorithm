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

func (l *LogEntryRepositoryImpl) Append(ctx context.Context, command string) error {
	//TODO implement me
	panic("implement me")
}

func (l *LogEntryRepositoryImpl) Get(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (l *LogEntryRepositoryImpl) GetFromIndex(ctx context.Context, index int) error {
	//TODO implement me
	panic("implement me")
}

func (l *LogEntryRepositoryImpl) DeleteFromIndex(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
