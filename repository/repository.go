package repository

import (
	"context"
)

type StateMachineRepository interface {
	create(ctx context.Context, key string, value int) error
	get(ctx context.Context, key string) error
	update(ctx context.Context, key string, value int) error
	delete(ctx context.Context, key string) error
}

type RaftLogRepository interface {
	append(ctx context.Context, index int) error
	get(ctx context.Context) error
	getFromIndex(ctx context.Context) error
	deleteFromIndex(ctx context.Context) error
}
