package web

import (
	"context"
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

	cleanup, err := ctxval.Resolve[ctxval.CleanupRecorder](ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving cleanup recorder: %w", err)
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
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			applog.Error(ctx, "HTTP server ListenAndServe error", "error", err)
		}
	}()

	cleanup(closeServer(server))

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

	router.Use(
		injectApplicationContext(ctx),       // must be first, since other middlewares may rely on it for dependencies
		middleware.RequestID,                // assign a request ID to the request
		requestIDStructuredLog,              // ensure the request ID is part of every log entry
		middleware.Logger,                   // log requests
		middleware.Timeout(timeoutDuration), // set a timeout for requests
	)

	strictHandler := api.NewStrictHandler(controller, middlewares)
	handler := api.HandlerFromMux(strictHandler, router)

	return handler, nil
}

func closeServer(server *http.Server) ctxval.CleanupFunc {
	return ctxval.CleanupFunc(func(ctx context.Context) {
		applog.Info(ctx, "Shutting down HTTP server...")
		if err := server.Shutdown(context.Background()); err == nil {
			applog.Info(ctx, "Shutting down HTTP server completed")
		} else {
			applog.Error(ctx, "HTTP server Shutdown error", "error", err)
		}
	})
}
