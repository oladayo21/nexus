package api

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type APIServer struct {
	e    *echo.Echo
	opts *APIServerOptions
}

type APIServerOptions struct {
	Port        int
	ServeStatic bool
	WebFS       fs.FS
}

func NewAPIServer(opts *APIServerOptions) *APIServer {
	e := echo.New()
	e.HideBanner = true

	registerRoutes(e)

	if opts.ServeStatic && opts.WebFS != nil {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       ".",
			Index:      "index.html",
			HTML5:      true,
			Filesystem: http.FS(opts.WebFS),
		}))
	}

	return &APIServer{
		e:    e,
		opts: opts,
	}
}

func registerRoutes(e *echo.Echo) {
	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})
}

func (s *APIServer) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.opts.Port))
}

func (s *APIServer) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
