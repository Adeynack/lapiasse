package controller

import (
	"fmt"
	"net/http"

	"adeynack.net/lapiasse/pkg/api"
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

func api422Error(err error) api.Error {
	return api.Error{
		Type:   api.ErrorTypeErrorValidation,
		Title:  "Resource did not validate",
		Detail: lo.ToPtr(err.Error()),
		Status: http.StatusUnprocessableEntity,
	}
}
