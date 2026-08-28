package repository

import (
	"context"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type KeyValueRepository interface {
	CreateBucket(ctx context.Context)
	create(ctx context.Context, key string, value int) error
	get(ctx context.Context, key string) error
	update(ctx context.Context, key string, value int) error
	delete(ctx context.Context, key string) error
}

type KeyValueRepositoryImpl struct {
	db *bolt.DB
}

func (r *KeyValueRepositoryImpl) get(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func (r *KeyValueRepositoryImpl) update(ctx context.Context, key string, value int) error {
	//TODO implement me
	panic("implement me")
}

func (r *KeyValueRepositoryImpl) delete(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func (r *KeyValueRepositoryImpl) CreateBucket(ctx context.Context) {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("KeyValue"))
		_ = b
		// TODO maybe use b
		if err != nil {
			return fmt.Errorf("Create Bucket has failed %s", err)
		}
		return nil
	})
	_ = err
	// TODO use error
}

func (r *KeyValueRepositoryImpl) create(ctx context.Context, key string, value int) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("KeyValue"))
		err := b.Put([]byte(key), []byte(string(value)))
		if err != nil {
			return fmt.Errorf("Create Bucket has failed %s", err)
		}
		return nil
	})
}

func (r *KeyValueRepositoryImpl) save(ctx context.Context) {
	err := r.db.Batch(func(tx *bolt.Tx) error {
		return nil
	})
	_ = err
	// TODO use error
}
