//go:build test

package validatorex

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func RequireValidationErrorTags(t *testing.T, err error, expected map[string][]string) {
	t.Helper()

	var validationErr validator.ValidationErrors
	require.ErrorAs(t, err, &validationErr)

	// Map of field name to map of tag name to FieldError
	actual := make(map[string]map[string]validator.FieldError)
	for _, fe := range validationErr {
		errorsForField, ok := actual[fe.Field()]
		if !ok {
			errorsForField = make(map[string]validator.FieldError)
			actual[fe.Field()] = errorsForField
		}

		errorsForField[fe.Tag()] = fe
	}

	var errors strings.Builder

	for field, expectedTags := range expected {
		actualTags, ok := actual[field]
		if !ok {
			errors.WriteString(fmt.Sprintf("Expected field %q not found in validation errors\n", field))
			continue
		}

		for _, tag := range expectedTags {
			if _, ok := actualTags[tag]; !ok {
				errors.WriteString(fmt.Sprintf("Expected validation error for field %q with tag %q not found\n", field, tag))
				continue
			}

			delete(actualTags, tag)
		}

		for tag := range actualTags {
			errors.WriteString(fmt.Sprintf("Unexpected validation error for field %q with tag %q\n", field, tag))
		}

		delete(actual, field)
	}

	for field, actualTags := range actual {
		for tag := range actualTags {
			errors.WriteString(fmt.Sprintf("Unexpected validation error for field %q with tag %q\n", field, tag))
		}
	}

	if errors.Len() > 0 {
		require.Failf(t, "Validation errors did not match expected", errors.String())
	}
}

func RequireValidationError(t *testing.T, err error, field string, tag string) {
	t.Helper()

	var validationErr validator.ValidationErrors
	require.ErrorAs(t, err, &validationErr)

	for _, fe := range validationErr {
		if fe.Field() == field && fe.Tag() == tag {
			return
		}
	}

	require.Failf(t, "Expected validation error not found", "field: %s, tag: %s\nerrors:\n%s", field, tag, validationErr)
}
