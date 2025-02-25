package corefx

import (
	"log/slog"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/lib/slogutil"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"core",

	// logger
	fx.Invoke(func(cfg *env.Config) {
		slog.SetLogLoggerLevel(slogutil.LogLevel(cfg.Log.Level))
	}),
)
