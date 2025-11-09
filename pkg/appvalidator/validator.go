package appvalidator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	defaultInstanceOnce sync.Once
	defaultInstance     *validator.Validate
)

func Default() *validator.Validate {
	defaultInstanceOnce.Do(initializeDefaultInstance)

	return defaultInstance
}

func initializeDefaultInstance() {
	defaultInstance = validator.New(
		validator.WithRequiredStructEnabled(),
	)
}
