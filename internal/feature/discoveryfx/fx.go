package discoveryfx

import (
	"context"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/feature/discovery"
	"github.com/inhibitor1217/moru/internal/lib/beacon"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"discovery",

	fx.Provide(func(cfg *env.Config) beacon.UDPBroadcastConfig {
		return beacon.UDPBroadcastConfig{
			Port: cfg.Discovery.Port,
		}
	}),
	fx.Provide(beacon.NewUDPBroadcast),
	fx.Provide(discovery.NewLocalDiscoverySvc),

	fx.Invoke(func(lc fx.Lifecycle, s discovery.DiscoverySvc) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.Start(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return s.Stop(ctx)
			},
		})
	}),

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
