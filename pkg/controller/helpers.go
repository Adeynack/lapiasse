package controller

import (
	"adeynack.net/lapiasse/pkg/api"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 50 // if this changes, change the maximum in api.openapi.yaml as well for parameter "page_size".
)

func scopePaginate(
	page *api.Page,
	pageSize *api.PageSize,
) func(tx *gorm.Statement) {
	limit := max(1, min(lo.FromPtr(pageSize), DefaultPageSize))

	pageNumber := max(lo.FromPtr(page), 1)

	offset := (pageNumber - 1) * limit

	return func(db *gorm.Statement) {
		db.Offset(offset).Limit(limit)
	}
}
