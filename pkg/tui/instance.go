package tui

import (
	"errors"
	"fmt"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/app"
	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
)

type Instance struct {
	App *app.Instance
}

func NewInstance(appInstance *app.Instance) (*Instance, error) {
	if appInstance == nil {
		return nil, errors.New("app instance is nil")
	}

	return &Instance{
		App: appInstance,
	}, nil
}

func (i *Instance) Close() error {
	if i == nil {
		return nil
	}

	return i.App.Close()
}

func (i *Instance) Run() error {
	ctx := i.App.Context()
	applog.Info(i.App.Context(), "Running TUI")

	// Placeholder TUI
	repo := ctxval.MustResolve[api.StrictServerInterface](ctx)

	booksResponse, err := repo.GetBooks(ctx, api.GetBooksRequestObject{})
	if err != nil {
		return fmt.Errorf("failed to get books: %w", err)
	}
	switch resp := booksResponse.(type) {
	case api.GetBooks200JSONResponse:
		fmt.Println("Books:")
		for i, book := range resp.Books {
			fmt.Printf("%d: %s\n", i+1, book.Name)
		}
	default:
		return fmt.Errorf("unexpected response type: %T", booksResponse)
	}

	fmt.Println("Press ENTER to exit")
	_, _ = fmt.Scanln()

	return nil
}
