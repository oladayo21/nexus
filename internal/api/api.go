package api

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/oladayo21/nexus/internal/config"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	Config *config.Config
	Redis  *redis.Client
	WebFS  fs.FS
}

type APIServer struct {
	echo           *echo.Echo
	opts           *Options
	sessionManager *scs.SessionManager
}

func NewAPIServer(opts *Options) *APIServer {
	sm := scs.New()
	sm.Store = goredisstore.New(opts.Redis)
	sm.Lifetime = 7 * 24 * time.Hour
	sm.IdleTimeout = 24 * time.Hour
	sm.Cookie.Name = "session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.Secure = opts.Config.IsProduction()
	sm.Cookie.SameSite = http.SameSiteLaxMode

	s := &APIServer{
		echo:           echo.New(),
		opts:           opts,
		sessionManager: sm,
	}

	s.echo.HideBanner = true
	s.registerMiddlewares()
	s.registerRoutes()

	return s
}

func (s *APIServer) Start() error {
	return s.echo.Start(fmt.Sprintf(":%d", s.opts.Config.Port))
}

func (s *APIServer) Stop(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
