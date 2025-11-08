package web

import (
	"context"
	"net/http"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// injectApplicationContext creates a middleware that injects the application context
// into the request context as a fallback for resolving values. The request's context
// values take precedence over the application context values.
func injectApplicationContext(appContext context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set the application context as fallback for resolving values.
			valuesJoined := ctxval.FallbackValues(r.Context(), appContext)

			next.ServeHTTP(w, r.WithContext(valuesJoined))
		})
	}
}

// requestIDStructuredLog is a middleware that adds the request ID to the
// structured log context.
func requestIDStructuredLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestId := ctx.Value(middleware.RequestIDKey)
		ctx = applog.With(ctx, "request_id", requestId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// logRequest creates a middleware that logs HTTP requests using a custom log formatter.
func logRequest() func(http.Handler) http.Handler {
	formatter := &logFormatter{
		SplitLogs: false,
	}

	return middleware.RequestLogger(formatter)
}

// handleCors creates a middleware that handles CORS (Cross-Origin Resource Sharing).
func handleCors(ctx context.Context) func(http.Handler) http.Handler {
	if env.GetRunEnv(ctx) == env.EnvProduction {
		// In production, CORS should be handled by the reverse proxy (e.g. Nginx).
		return nil
	}

	// Basic CORS example taken from https://github.com/go-chi/cors#usage
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
