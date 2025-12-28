package web

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"runtime/debug"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/samber/lo"
)

var internalServerJsonError = api.Error{
	Status: http.StatusInternalServerError,
	Title:  "Internal Server Error",
	Type:   api.ErrorTypeErrorInternalError,
}

var generic500ErrorJson = lo.Must(json.Marshal(internalServerJsonError))

func respondWithJsonError(w http.ResponseWriter, r *http.Request, jsonError api.Error) {
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(jsonError)
	if err != nil {
		applog.Error(r.Context(), "Failed to marshal error response", "error", err)
		http.Error(w, string(generic500ErrorJson), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(jsonError.Status)
	if _, err := w.Write(json); err != nil {
		applog.Error(r.Context(), "Failed to write error response", "error", err)
	}
}

func requestErrorHandler() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		applog.Debug(r.Context(), "Failed to handle request", "error", err)
		respondWithJsonError(w, r, api.Error{
			Detail: lo.ToPtr(err.Error()),
			Status: http.StatusBadRequest,
			Title:  "Bad Request",
			Type:   api.ErrorTypeErrorBadRequest,
		})
	}
}

func responseErrorHandler() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		applog.Error(r.Context(), "Failed to handle response", "error", err)
		respondWithJsonError(w, r, internalServerJsonError)
	}
}

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

func recoverFromPanicAsJsonErr() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Inspired by [middleware.Recoverer]
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rvr)
					}

					applog.Error(r.Context(), "Panic recovered in HTTP request", "error", rvr)

					if logEntry := middleware.GetLogEntry(r); logEntry != nil {
						logEntry.Panic(rvr, debug.Stack())
					}

					if r.Header.Get("Connection") != "Upgrade" {
						respondWithJsonError(w, r, internalServerJsonError)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithJsonError(w, r, api.Error{
			Status: http.StatusNotFound,
			Title:  "Not Found",
			Type:   api.ErrorTypeErrorNotFound,
			Detail: lo.ToPtr(fmt.Sprintf("Path %q not found", r.URL.Path)),
		})
	}
}

func handleMethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithJsonError(w, r, api.Error{
			Status: http.StatusMethodNotAllowed,
			Title:  "Method Not Allowed",
			Type:   api.ErrorTypeErrorMethodNotAllowed,
			Detail: lo.ToPtr(fmt.Sprintf("Path %q does not allow verb %q", r.URL.Path, r.Method)),
		})
	}
}
