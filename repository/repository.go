package repository

import (
	"context"
)

type DbRepository interface {
	create(ctx context.Context, key string, value int) error
	update(ctx context.Context, key string, value int) error
	delete(ctx context.Context, key string) error
}
