package repository

import (
	"context"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type LogEntryRepository interface {
	CreateLogBucket(ctx context.Context) error
	Append(ctx context.Context, command string, term int) error
	Get(ctx context.Context, command string) (int, error)
	GetFromIndex(ctx context.Context, index int) error
	DeleteFromIndex(ctx context.Context) error
}

type LogEntryRepositoryImpl struct {
	db *bolt.DB
}

func (l *LogEntryRepositoryImpl) CreateLogBucket(ctx context.Context) error {
	return l.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("logs"))
		return err
	})
}

func (l *LogEntryRepositoryImpl) Append(ctx context.Context, command string, term int) error {
	return l.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("logs"))
		if b == nil {
			return fmt.Errorf("logs bucket does not exist")
		}

		return b.Put([]byte(command), []byte{byte(term)})
	})
}

func (l *LogEntryRepositoryImpl) Get(ctx context.Context, command string) (int, error) {
	var term int

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("logs"))
		if b == nil {
			return fmt.Errorf("logs bucket does not exist")
		}

		value := b.Get([]byte(command))
		if value == nil {
			return fmt.Errorf("log entry not found")
		}

		term = int(value[0])

		return nil
	})

	return term, err
}

func (l *LogEntryRepositoryImpl) GetFromIndex(ctx context.Context, index int) error {
	//TODO implement me
	panic("implement me")
}

func (l *LogEntryRepositoryImpl) DeleteFromIndex(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
