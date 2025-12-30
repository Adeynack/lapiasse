package model

import (
	"context"
	"fmt"
	"strings"

	"adeynack.net/lapiasse/pkg/api"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
)

type ValidationError struct {
	FieldErrors []FieldError
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %v", e.FieldErrors)
}

type ModelValidatable interface {
	Validate(ctx context.Context, reason ValidationReason) error
}

type FieldError api.FieldValidationFailure

type ValidationErrorBuilder []FieldError

func (b *ValidationErrorBuilder) ToError() error {
	if len(*b) == 0 {
		return nil
	}

	return ValidationError{FieldErrors: *b}
}

func (b *ValidationErrorBuilder) Add(field, message, validation string, param string) {
	*b = append(*b, FieldError{
		Field:      field,
		Message:    message,
		Validation: validation,
		Param:      lo.EmptyableToPtr(param),
	})
}

func (b *ValidationErrorBuilder) AddFromValidator(e validator.ValidationErrors) {
	for _, ve := range e {
		b.Add(
			ve.Namespace(),
			validationMessageFromFieldError(ve),
			ve.ActualTag(),
			ve.Param(),
		)
	}
}

func validationMessageFromFieldError(fe validator.FieldError) string {
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
