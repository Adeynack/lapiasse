//go:build test

package repository

import (
	"context"

	"github.com/samber/lo"
)

func MustCreate[T any](ctx context.Context, entity *T) *T {
	return lo.Must1(Create(ctx, entity))
}

func MustCreate0[T any](ctx context.Context, entity *T) {
	_ = lo.Must1(Create(ctx, entity))
}
