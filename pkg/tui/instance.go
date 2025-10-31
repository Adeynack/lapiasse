package tui

import (
	"errors"
	"log/slog"

	"adeynack.net/lapiasse/pkg/app"
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
	slog.Info("Starting TUI")

	return nil
}
