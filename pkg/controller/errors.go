package controller

import (
	"fmt"

	"adeynack.net/lapiasse/pkg/api"
	"github.com/samber/lo"
)

func api404Error(instance string, detail string) api.Error {
	return api.Error{
		Status:   404,
		Title:    "Not Found",
		Detail:   &detail,
		Instance: &instance,
		Type:     lo.ToPtr("https://adeynack.net/lapiasse/errors/NotFound"),
	}
}

func api404ErrorFromId(
	resourceName string, // e.g. "Book"
	resourcePath string, // e.g. "/books"
	id string,
) api.Error {
	return api404Error(
		fmt.Sprintf("%s/%s", resourcePath, id),
		fmt.Sprintf("%s with ID %q not found", resourceName, id),
	)
}
