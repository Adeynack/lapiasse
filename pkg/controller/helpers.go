package controller

import (
	"adeynack.net/lapiasse/pkg/api"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 50    // if this changes, change the maximum in api.openapi.yaml as well for parameter "page_size".
	MaxPageSize     = 1_000 // if this changes, change the maximum in api.openapi.yaml as well for parameter "page_size".
)

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
