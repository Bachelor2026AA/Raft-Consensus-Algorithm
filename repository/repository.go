package repository

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"raft-consensus/domain"

	bolt "go.etcd.io/bbolt"
)

type LogEntryRepository interface {
	Append(entry domain.Log) error
	Get(index uint64) (domain.Log, error)
	GetFromIndex(ctx context.Context, index int) error
	DeleteFromIndex(ctx context.Context) error
}

type LogEntryRepositoryImpl struct {
	db *bolt.DB
}

func NewLogEntryRepository(db *bolt.DB) (LogEntryRepository, error) {
	repo := &LogEntryRepositoryImpl{
		db: db}
	if err := repo.createLogBucket(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (l *LogEntryRepositoryImpl) createLogBucket() error {
	return l.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("log_entries"))
		return err
	})
}

func (l *LogEntryRepositoryImpl) Append(entry domain.Log) error {
	return l.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("logs"))
		if b == nil {
			return fmt.Errorf("logs bucket does not exist")
		}

		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, entry.Index)

		return b.Put(key, data)
	})
}

func (l *LogEntryRepositoryImpl) Get(index uint64) (domain.Log, error) {
	var entry domain.Log

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("logs"))
		if b == nil {
			return fmt.Errorf("logs bucket does not exist")
		}

		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, index)

		value := b.Get(key)
		if value == nil {
			return fmt.Errorf("log entry not found")
		}

		if err := json.Unmarshal(value, &entry); err != nil {
			return fmt.Errorf("failed to decode log entry: %w", err)
		}

		return nil
	})

	return entry, err
}

func (l *LogEntryRepositoryImpl) GetFromIndex(ctx context.Context, index int) error {
	//TODO implement me
	panic("implement me")
}

func (l *LogEntryRepositoryImpl) DeleteFromIndex(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
