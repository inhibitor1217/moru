package main

// #include <stdlib.h>
//
// typedef void (*log_write_t)(const void* msg, int len);
// void bridge_log_write(log_write_t f, const void* msg, int len);
//
// typedef struct {
//     const void* data;
//     int len;
// } ffi_t;
import "C"

import (
	"context"
	"log/slog"
	"sync"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/envfx"
	"github.com/inhibitor1217/moru/internal/feature/corefx"
	"github.com/inhibitor1217/moru/internal/feature/discovery"
	"github.com/inhibitor1217/moru/internal/feature/discoveryfx"
	"github.com/inhibitor1217/moru/internal/lib/slogutil"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type moru struct {
	mu sync.Mutex

	app     *fx.App
	logger  *nativeLogger
	started bool

	discoverySvc discovery.DiscoverySvc
}

var m = &moru{}

//export moru_init
func moru_init() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.app != nil {
		return
	}

	m.logger = newNativeLogger()

	m.app = fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),

		fx.Invoke(func(cfg *env.Config) {
			slog.SetDefault(slog.New(slog.NewTextHandler(m.logger, &slog.HandlerOptions{
				Level: slogutil.LogLevel(cfg.Log.Level),
			})))
		}),

		envfx.Option,

		corefx.Module,
		discoveryfx.Module,

		fx.Invoke(func(lc fx.Lifecycle, cfg *env.Config, discoverySvc discovery.DiscoverySvc) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logStart(cfg)
					return nil
				},
			})
		}),

		fx.Populate(&m.discoverySvc),
	)
}

//export moru_run
func moru_run() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.app == nil {
		slog.Warn("moru not initialized. moru_init() must be called first")
		return
	}

	if m.started {
		slog.Warn("moru already started")
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

	m.logger.Bind(fn)
}

//export moru_known_peers
func moru_known_peers(req C.ffi_t) C.ffi_t {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := context.Background()
	log := slog.Default().With("entrypoint", "moru_known_peers")

	if m.app == nil {
		log.Warn("moru not initialized. moru_init() must be called first")
		return C.ffi_t{}
	}

	if req.data == nil {
		log.Warn("invalid request: empty data buffer")
		return C.ffi_t{}
	}

	res := m.knownPeers(ctx, log, C.GoBytes(req.data, C.int(req.len)))

	if res == nil {
		log.Warn("failed to process request")
		return C.ffi_t{}
	}

	return C.ffi_t{data: C.CBytes(res), len: C.int(len(res))}
}

func logStart(cfg *env.Config) {
	slog.Default().
		WithGroup("application").
		With("name", cfg.Application.Name).
		With("stage", cfg.Application.Stage).
		Info("initializing moru")
}

func main() {}
