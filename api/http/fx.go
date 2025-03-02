package http

import (
	"github.com/inhibitor1217/moru/api/http/routes/ping"
	"github.com/inhibitor1217/moru/internal/lib/http"
	"go.uber.org/fx"
)

var Option = fx.Options(
	fx.Provide(asRoute(ping.New)),
)

func asRoute(route interface{}) interface{} {
	return fx.Annotate(route, fx.As(new(http.Routes)), fx.ResultTags(`group:"api.http.routes"`))
}
