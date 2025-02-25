package envfx

import (
	"github.com/inhibitor1217/moru/internal/env"
	"go.uber.org/fx"
)

var Option = fx.Options(
	fx.Provide(env.LoadViper),
)
