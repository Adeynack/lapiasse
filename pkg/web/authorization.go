package web

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/env"
)

func generateSessionToken() string {
	return rand.Text() + rand.Text()
}

const (
	StaticApiTokenDevelopment = "dev-static-token"
	StaticApiTokenTest        = "test-static-token"
)

// apiTokenMiddleware checks for the presence of a valid API token in the
// Authorization header of incoming requests.
//
// This is an interim solution in order to have some basic level of security
// until a proper authentication mechanism (e.g.: OAuth2, JWT) is implemented.
func apiTokenMiddleware(ctx context.Context) func(http.Handler) http.Handler {
	var expectedToken string
	switch env.GetRunEnv(ctx) {
	case env.EnvDevelopment:
		expectedToken = StaticApiTokenDevelopment
		applog.Warn(ctx, "Using static API token for development", "api_token", expectedToken)
	case env.EnvTest:
		expectedToken = StaticApiTokenTest
		applog.Warn(ctx, "Using static API token for test", "api_token", expectedToken)
	case env.EnvProduction:
		expectedToken = generateSessionToken()
		applog.Info(ctx, "Generated API token for production environment", "api_token", expectedToken)
	default:
		panic("unknown run environment")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestToken := r.Header.Get("Authorization")
			if requestToken != fmt.Sprintf("Bearer %s", expectedToken) {
				applog.Warn(r.Context(), "Unauthorized API request", "remote_addr", r.RemoteAddr)
				respondWithJsonError(w, r, api.Error{
					Status: http.StatusUnauthorized,
					Title:  "Unauthorized",
					Type:   api.ErrorTypeErrorUnauthorized,
				})

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
