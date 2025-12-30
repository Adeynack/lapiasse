package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/model"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 50    // if this changes, change the maximum in api.openapi.yaml as well for parameter "page_size".
	MaxPageSize     = 1_000 // if this changes, change the maximum in api.openapi.yaml as well for parameter "page_size".
)

// scopePaginate applies pagination to a GORM query based on API pagination parameters.
func scopePaginate(
	page *api.ReqPaginatedPage,
	pageSize *api.ReqPaginatedPageSize,
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

// validate ensures that a model is validated for a given reason.
func validate(ctx context.Context, value model.ModelValidatable, reason model.ValidationReason) validationResult {
	err := value.Validate(ctx, reason)

	if err == nil {
		return validationResult{}
	}

	var validationErr model.ValidationError
	if errors.As(err, &validationErr) {
		return validationResult{apiError: api422Error(validationErr)}
	}

	modelName := strings.TrimPrefix(fmt.Sprintf("%T", value), "*model.")
	return validationResult{err: fmt.Errorf("validating %s", modelName)}
}

type validationResult struct {
	apiError *api.ValidationError
	err      error
}

func (v validationResult) IsError() bool {
	return v.apiError != nil || v.err != nil
}

func (v validationResult) Unwrap() (api.ValidationError, error) {
	return lo.FromPtr(v.apiError), v.err
}
