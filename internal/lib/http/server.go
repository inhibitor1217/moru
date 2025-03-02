package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Router interface {
	gin.IRouter
}

type Routes interface {
	Path() string
	Register(Router)
}

type Middleware interface {
	Handler() gin.HandlerFunc
}

type Server struct {
	server *http.Server
	engine *gin.Engine
	params ServerParams
}

type ServerParams struct {
	Debug bool
	Port  int
}

func NewServer(
	params ServerParams,
	routes []Routes,
	middlewares []Middleware,
) (*Server, error) {
	engine := gin.Default()
	if params.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &Server{
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(params.Port),
			Handler: engine,
		},
		engine: engine,
		params: params,
	}

	for _, middleware := range middlewares {
		engine.Use(middleware.Handler())
	}

	for _, route := range routes {
		route.Register(engine.Group(route.Path()))
	}

	return server, nil
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
