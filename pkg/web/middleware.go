package web

import (
	"context"
	"net/http"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"github.com/go-chi/chi/v5/middleware"
)

// injectApplicationContext is a middleware that injects the application context
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
