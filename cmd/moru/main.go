package main

// #include <stdlib.h>
//
// typedef void (*log_write_t)(const void* msg, int len);
// void bridge_log_write(log_write_t f, const void* msg, int len);
import "C"

import (
	"context"
	"log/slog"
	"sync"
	"unsafe"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/envfx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type moru struct {
	mu sync.Mutex

	app     *fx.App
	log     C.log_write_t
	started bool
}

func (m *moru) Write(bs []byte) (int, error) {
	if m.log == nil {
		return len(bs), nil
	}

	if len(bs) == 0 {
		return 0, nil
	}

	p := C.CBytes(bs)
	// we cannot free the memory when this function returns, since
	// the memory is still being used asynchronously by the native code.

	C.bridge_log_write(m.log, p, C.int(len(bs)))
	return len(bs), nil
}

var m = &moru{}

//export moru_init
func moru_init() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.app != nil {
		return
	}

	m.app = fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),

		envfx.Option,

		fx.Invoke(func(lc fx.Lifecycle, cfg *env.Config) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logStart(cfg)
					return nil
				},
			})
		}),
	)

	slog.SetDefault(slog.New(slog.NewTextHandler(m, nil)))
}

//export moru_run
func moru_run() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.app == nil {
		return
	}

	if m.started {
		return
	}

	m.started = true

	// do not block the main goroutine in native bindings.
	// let's worry with graceful shutdown later.
	go m.app.Run()
}

//export moru_register_logger
func moru_register_logger(fn C.log_write_t) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log = fn
}

//export moru_log_ack
func moru_log_ack(msg unsafe.Pointer) {
	C.free(msg)
}

func logStart(cfg *env.Config) {
	slog.Default().
		WithGroup("application").
		With("name", cfg.Application.Name).
		With("stage", cfg.Application.Stage).
		Info("initializing moru")
}

func main() {}
