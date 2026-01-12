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

type ValidationBuilder struct {
	reason    ValidationReason
	namespace string
	errors    []FieldError
}

// Creating returns true if the validation reason is Creation.
func (b *ValidationBuilder) Creating() bool {
	return b.reason == ValidationReasonCreate
}

// Updating returns true if the validation reason is Update.
func (b *ValidationBuilder) Updating() bool {
	return b.reason == ValidationReasonUpdate
}

// Deleting returns true if the validation reason is Delete.
func (b *ValidationBuilder) Deleting() bool {
	return b.reason == ValidationReasonDelete
}

// ReasonIs checks if the validation reason matches any of the provided reasons.
func (b *ValidationBuilder) ReasonIs(reasons ...ValidationReason) bool {
	return lo.Contains(reasons, b.reason)
}

func (b *ValidationBuilder) Reason() ValidationReason {
	return b.reason
}

func (b *ValidationBuilder) ToError() error {
	if len(b.errors) == 0 {
		return nil
	}

	return ValidationError{FieldErrors: b.errors}
}

// Namespaced returns a new ValidationErrorBuilder with the given sub-namespace applied.
func (b *ValidationBuilder) Namespaced(namespace string) *ValidationBuilder {
	ns := namespace
	if b.namespace != "" {
		ns = b.namespace + "." + namespace
	}

	return &ValidationBuilder{
		namespace: ns,
		errors:    b.errors,
	}
}

// AddFieldErr adds a new field error to the builder.
func (b *ValidationBuilder) AddFieldErr(field, message, validation string, param string) {
	namespacedField := field
	if b.namespace != "" {
		namespacedField = b.namespace + "." + field
	}

	b.errors = append(b.errors, FieldError{
		Field:      namespacedField,
		Message:    message,
		Validation: validation,
		Param:      lo.EmptyableToPtr(param),
	})
}

// AddFromValidator adds field errors from a [validator.ValidationErrors] instance.
func (b *ValidationBuilder) AddFromValidator(e validator.ValidationErrors) {
	for _, ve := range e {
		ns := ve.Namespace()
		if nsParts := strings.Split(ns, "."); len(nsParts) > 0 && nsParts[0] == b.namespace {
			// Remove the prefix namespace if it matches the builder's namespace.
			ns = strings.Join(nsParts[1:], ".")
		}

		b.AddFieldErr(
			ns, // TODO: Test if deep namespace works correctly.
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
