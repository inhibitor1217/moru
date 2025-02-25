package discoveryfx

import (
	"context"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/lib/beacon"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"discovery",

	// UDP broadcast beacon
	fx.Provide(func(cfg *env.Config) beacon.UDPBroadcastConfig {
		return beacon.UDPBroadcastConfig{
			Port: cfg.Discovery.Port,
		}
	}),
	fx.Provide(beacon.NewUDPBroadcast),
	fx.Invoke(func(lc fx.Lifecycle, b beacon.Beacon) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return b.Start(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return b.Stop(ctx)
			},
		})
	}),
)
