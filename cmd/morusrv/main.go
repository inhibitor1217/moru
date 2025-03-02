package main

import (
	"log/slog"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/envfx"
	"github.com/inhibitor1217/moru/internal/feature/apifx"
	"github.com/inhibitor1217/moru/internal/feature/corefx"
	"github.com/inhibitor1217/moru/internal/feature/discoveryfx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),

		envfx.Option,

		apifx.Module,
		corefx.Module,
		discoveryfx.Module,

		fx.Invoke(logStart),
	).Run()
}

func logStart(cfg *env.Config) {
	slog.Default().
		WithGroup("application").
		With("name", cfg.Application.Name).
		With("role", cfg.Application.Role).
		With("stage", cfg.Application.Stage).
		Info("running")
}
