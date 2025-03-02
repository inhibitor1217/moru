package ping

import (
	netHTTP "net/http"

	"github.com/gin-gonic/gin"
	"github.com/inhibitor1217/moru/internal/lib/http"
)

type ping struct{}

func New() *ping {
	return &ping{}
}

func (p *ping) Path() string {
	return "/ping"
}

func (p *ping) Register(router http.Router) {
	router.GET("", p.Ping)
}

func (p *ping) Ping(ctx *gin.Context) {
	ctx.String(netHTTP.StatusOK, "pong")
}
