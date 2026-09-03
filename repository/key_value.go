package repository

import (
	"context"

	bolt "go.etcd.io/bbolt"
)

type KeyValueRepository interface {
	create(ctx context.Context, key string, value int) error
	get(ctx context.Context, key string) error
	update(ctx context.Context, key string, value int) error
	delete(ctx context.Context, key string) error
}

type KeyValueRepositoryImpl struct {
	db *bolt.DB
}

func NewKeyValueRepository(db *bolt.DB) (KeyValueRepository, error) {
	repo := &KeyValueRepositoryImpl{db: db}

	if err := repo.createKeyValueBucket(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *KeyValueRepositoryImpl) createKeyValueBucket() error {
	return r.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("KeyValue"))
		return err
	})
}
func (r *KeyValueRepositoryImpl) create(ctx context.Context, key string, value int) error {
	panic("implement me")
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
