package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"

	"adeynack.net/lapiasse/pkg/controller"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/repository"
	"adeynack.net/lapiasse/pkg/web"
)

// Instance represents a running instance of the application.
type Instance struct {
	dependenciesContext *ctxval.Container
	cancel              context.CancelFunc
	cleanupFuncs        []ctxval.CleanupFunc
	closeOnce           func() error
}

func NewInstance(ch *ConfigurationHolder) (i *Instance, err error) {
	if ch == nil {
		return nil, errors.New("configuration holder is nil")
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	container := ctxval.NewContainer(cancelCtx)
	i = &Instance{
		dependenciesContext: container,
		cancel:              cancel,
		cleanupFuncs:        make([]ctxval.CleanupFunc, 0),
	}
	i.closeOnce = sync.OnceValue(i.doCloseOnce)

	defer func() {
		if err != nil {
			i.Close()
			i = nil
		}
	}()

	// Register simple dependencies (no init, no error).
	ctxval.RegisterInContainer[ctxval.CleanupRecorder](container, func(f ctxval.CleanupFunc) {
		i.cleanupFuncs = append(i.cleanupFuncs, f)
	})
	ctxval.RegisterInContainer(container, slog.Default()) // temporary as default logger, until `configureLogger` runs
	ctxval.RegisterInContainer(container, ch)
	ctxval.RegisterInContainer(container, controller.New())

	// Register dependencies requiring initialization.
	c := ch.Configuration
	reg(i, &err, c.Data, repository.InitializeDataFilesystem, "initializing data folder")
	reg(i, &err, c.Log, configureLogger, "configuring logger")
	reg(i, &err, c.Data, repository.InitializeGorm, "initializing Gorm database")
	reg(i, &err, c.Web, web.StartServer, "starting web server")

	if err == nil {
		i.closeOnOsSignal()
	}

	return
}

// reg is a helper function to register a dependency in the instance's container.
func reg[T, P any](
	instance *Instance,
	err *error,
	param P,
	factory func(ctx context.Context, param P) (T, error),
	errorContext string,
) {
	if *err != nil {
		return
	}

	value, factoryErr := factory(instance.dependenciesContext, param)
	if factoryErr != nil {
		*err = fmt.Errorf("%s: %w", errorContext, factoryErr)

		return
	}

	ctxval.RegisterInContainer(instance.dependenciesContext, value)
}

func (instance *Instance) Context() context.Context {
	return instance.dependenciesContext
}

// Close implements the [io.Closer] interface.
func (instance *Instance) Close() error {
	return instance.closeOnce()
}

func (instance *Instance) doCloseOnce() error {
	if instance == nil {
		return nil
	}

	// Cancel the instance context to stop all other background operations.
	if instance.cancel != nil {
		instance.cancel()
	}

	// Call all registered cleanup functions, in reverse order.
	var errs []error
	for i := len(instance.cleanupFuncs) - 1; i >= 0; i-- {
		cleanupFunc := instance.cleanupFuncs[i]
		if cleanupFunc != nil {
			if err := cleanupFunc(instance.dependenciesContext); err != nil {
				// Do NOT early-return on error.
				// All dependencies must have a chance to clean up, and the context to close.
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)

	return nil // todo
}

func (instance *Instance) closeOnOsSignal() {
	go func() {
		interruptCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		<-interruptCtx.Done()

		_ = instance.Close()
	}()
}
