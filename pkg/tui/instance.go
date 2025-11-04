package tui

import (
	"errors"
	"fmt"

	"adeynack.net/lapiasse/pkg/app"
	"adeynack.net/lapiasse/pkg/applog"
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
	applog.Info(i.App.Context(), "Running TUI")

	// Temporary "TUI" ;-)
	fmt.Println("Press ENTER to exit")
	_, _ = fmt.Scanln()

	return nil
}
