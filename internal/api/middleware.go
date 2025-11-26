package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *APIServer) registerMiddlewares() {
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(s.sessionMiddleware())

	if s.opts.Config.IsProduction() && s.opts.WebFS != nil {
		s.echo.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       ".",
			Index:      "index.html",
			HTML5:      true,
			Filesystem: http.FS(s.opts.WebFS),
		}))
	}
}

func (s *APIServer) sessionMiddleware() echo.MiddlewareFunc {
	return echo.WrapMiddleware(s.sessionManager.LoadAndSave)
}
