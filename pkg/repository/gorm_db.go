package repository

import (
	"context"

	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"gorm.io/gorm"
)

func Create[T any](ctx context.Context, entity *T) (*T, error) {
	return withGorm(ctx, func(db *gorm.DB) (*T, error) {
		err := gorm.G[T](db).Create(ctx, entity)

		return entity, err
	})
}

func withGorm[T any](ctx context.Context, fn func(db *gorm.DB) (*T, error)) (*T, error) {
	db, err := ctxval.Resolve[*gorm.DB](ctx)
	if err != nil {
		return nil, err
	}

	return fn(db)
}
