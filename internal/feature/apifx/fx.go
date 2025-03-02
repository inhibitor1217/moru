package apifx

import (
	"context"
	"log/slog"
	stdHTTP "net/http"
	"time"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/lib/http"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"api",

	fx.Provide(func(cfg *env.Config) http.ServerParams {
		return http.ServerParams{
			Port: cfg.HTTP.Port,
		}
	}),

	fx.Provide(fx.Annotate(
		http.NewServer,
		fx.ParamTags(
			``,
			`group:"api.http.routes"`,
			`group:"api.http.middlewares"`,
		),
		fx.OnStart(func(ctx context.Context, server *http.Server) error {
			go func() {
				slog.Default().InfoContext(ctx, "api HTTP server started")
				if err := server.Start(); err != nil && err != stdHTTP.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		}),
		fx.OnStop(func(ctx context.Context, server *http.Server) error {
			slog.Default().InfoContext(ctx, "api HTTP server stopped")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return server.Stop(shutdownCtx)
		}),
	)),

	// require *http.Server to be injected
	fx.Invoke(func(*http.Server) {}),
)
