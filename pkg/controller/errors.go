package controller

import (
	"fmt"
	"net/http"
	"strings"

	"adeynack.net/lapiasse/pkg/api"
	"github.com/go-playground/validator/v10"
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

func api422Error(errs validator.ValidationErrors) *api.ValidationError {
	validationErrors := lo.Map(errs, func(fe validator.FieldError, _ int) api.FieldValidationError {
		return api.FieldValidationError{
			Field:      fe.Namespace(),
			Message:    validationMessage(fe),
			Validation: fe.ActualTag(),
			Param:      lo.EmptyableToPtr(fe.Param()),
		}
	})

	return &api.ValidationError{
		Type:             api.ErrorTypeErrorValidation,
		Title:            "Resource did not validate",
		Status:           http.StatusUnprocessableEntity,
		ValidationErrors: validationErrors,
	}
}

func validationMessage(fe validator.FieldError) string {
	parts := make([]string, 0, 2)

	switch fe.Tag() {
	case "currencyIsoCode":
		parts = append(parts, "currency ISO code")
	}

	switch fe.ActualTag() {
	case "required":
		parts = append(parts, "is required")
	case "len":
		parts = append(parts, fmt.Sprintf("must be %s characters long", fe.Param()))
	default:
		parts = append(parts, fmt.Sprintf("failed validation %s(%s)", fe.ActualTag(), fe.Param()))
	}

	return strings.Join(parts, " ")
}
