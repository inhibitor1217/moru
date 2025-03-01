package main

import "C"

import (
	"log/slog"
	"sync"
	"time"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/envfx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type state struct {
	mu sync.Mutex

	app     *fx.App
	started bool
}

var s = &state{}

//export moru_init
func moru_init() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.app != nil {
		return
	}

	s.app = fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),

		envfx.Option,

		fx.Invoke(logStart),
		fx.Invoke(ping),
	)
}

//export moru_run
func moru_run() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.app == nil {
		return
	}

	if s.started {
		return
	}

	s.started = true

	// do not block the main goroutine in native bindings.
	// let's worry with graceful shutdown later.
	go s.app.Run()
}

func logStart(cfg *env.Config) {
	slog.Default().
		WithGroup("application").
		With("name", cfg.Application.Name).
		With("stage", cfg.Application.Stage).
		Info("initializing moru")
}

func ping() {
	tick := time.NewTicker(5 * time.Second)

	go func() {
		for range tick.C {
			slog.Default().Info("ping")
		}
	}()
}

func main() {}
