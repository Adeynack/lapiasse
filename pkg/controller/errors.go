package controller

import (
	"fmt"
	"net/http"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/model"
	"github.com/samber/lo"
)

func api404Error(detail string) api.Error {
	return api.Error{
		Status: 404,
		Title:  "Not Found",
		Detail: &detail,
		Type:   api.ErrorTypeErrorNotFound,
	}
}

func api404ErrorFromId(
	resourceName string, // e.g. "Book"
	id string,
) api.Error {
	return api404Error(
		fmt.Sprintf("%s with ID %q not found", resourceName, id),
	)
}

func api422Error(errs model.ValidationError) *api.ValidationError {
	validationErrors := lo.Map(errs.FieldErrors, func(e model.FieldError, _ int) api.FieldValidationError {
		return api.FieldValidationError(e)
	})

	return &api.ValidationError{
		Type:             api.ErrorTypeErrorValidation,
		Title:            "Resource did not validate",
		Status:           http.StatusUnprocessableEntity,
		ValidationErrors: validationErrors,
	}
}
