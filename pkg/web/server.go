package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartServer(ctx context.Context, config *Configuration) (*http.Server, error) {
	if !config.Expose {
		applog.Debug(ctx, "Web server is disabled by configuration")
		return &http.Server{}, nil
	}

	handler, err := createHandler(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP handler: %w", err)
	}

	address := fmt.Sprintf("localhost:%d", config.Port)
	server := &http.Server{Addr: address, Handler: handler}

	// Start the server in the background
	go func() {
		applog.Info(ctx, "Starting HTTP server", "address", address)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				applog.Info(ctx, "HTTP server closed")
			} else {
				applog.Error(ctx, "HTTP server ListenAndServe error", "error", err)
			}
		}
	}()

	ctxval.MustCleanup(ctx, closeServer(server))

	return server, nil
}

func createHandler(ctx context.Context, config *Configuration) (http.Handler, error) {
	controller, err := ctxval.Resolve[api.StrictServerInterface](ctx)
	if err != nil {
		return nil, err
	}

	middlewares := []api.StrictMiddlewareFunc{}
	router := chi.NewMux()
	timeoutDuration := time.Duration(config.RequestTimeoutMs) * time.Millisecond

	routerUseNonNilMiddlewares(router,
		injectApplicationContext(ctx),       // must be first, since other middlewares may rely on it for dependencies
		middleware.RequestID,                // assign a request ID to the request
		requestIDStructuredLog,              // ensure the request ID is part of every log entry
		logRequest(),                        // log requests
		recoverFromPanicAsJsonErr(),         // recover from panics
		handleCors(ctx),                     // handle CORS
		middleware.Timeout(timeoutDuration), // set a timeout for requests
		apiTokenMiddleware(ctx),             // check for API token in requests (for all paths)
	)

	router.NotFound(handleNotFound())
	router.MethodNotAllowed(handleMethodNotAllowed())

	strictHandler := api.NewStrictHandlerWithOptions(controller, middlewares, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler(),
		ResponseErrorHandlerFunc: responseErrorHandler(),
	})
	handler := api.HandlerFromMux(strictHandler, router)

	return handler, nil
}

func routerUseNonNilMiddlewares(router chi.Router, middlewares ...func(http.Handler) http.Handler) {
	for _, mw := range middlewares {
		if mw != nil {
			router.Use(mw)
		}
	}
}

func closeServer(server *http.Server) ctxval.CleanupFunc {
	return ctxval.CleanupFunc(func(ctx context.Context) error {
		applog.Info(ctx, "Shutting down HTTP server...")
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("shutting down HTTP server: %w", err)
		}

		return nil
	})
}
