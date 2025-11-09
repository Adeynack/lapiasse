package controller

import (
	"context"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/appvalidator"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 50    // if this changes, change the maximum in api.openapi.yaml as well for parameter "page_size".
	MaxPageSize     = 1_000 // if this changes, change the maximum in api.openapi.yaml as well for parameter "page_size".
)

// scopePaginate applies pagination to a GORM query based on API pagination parameters.
func scopePaginate(
	page *api.Page,
	pageSize *api.PageSize,
) func(tx *gorm.Statement) {
	limit := lo.FromPtrOr(pageSize, DefaultPageSize) // get or default
	limit = max(1, limit)                            // at least 1
	limit = min(limit, DefaultPageSize)              // at most DefaultPageSize

	pageNumber := lo.FromPtr(page)  // get or default to 0
	pageNumber = max(pageNumber, 1) // at least 1

	offset := (pageNumber - 1) * limit // database offset

	return func(db *gorm.Statement) {
		db.Offset(offset).Limit(limit)
	}
}

// validate ensures that a model is validated and otherwise returns
// a pre-filled [api.Error] containing the details of the validation fail.
func validate(ctx context.Context, value any) (bool, api.Error) {
	err := appvalidator.Default().StructCtx(ctx, value)
	if err == nil {
		return true, api.Error{}
	}

	return false, api422Error(err)
}
