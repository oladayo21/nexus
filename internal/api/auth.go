package api

import (
	"github.com/labstack/echo/v4"
)

func (s *APIServer) GetUserID(c echo.Context) (int64, bool) {
	userID := s.sessionManager.GetInt64(c.Request().Context(), "userID")

	return userID, userID != 0
}

func (s *APIServer) SetUserID(c echo.Context, userID int64) error {
	if err := s.sessionManager.RenewToken(c.Request().Context()); err != nil {
		return err
	}

	s.sessionManager.Put(c.Request().Context(), "userID", userID)

	return nil
}

func (s *APIServer) ClearSession(c echo.Context) error {
	return s.sessionManager.Destroy(c.Request().Context())
}

func (s *APIServer) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, ok := s.GetUserID(c); !ok {
			return echo.NewHTTPError(401, "unauthorized")
		}

		return next(c)
	}
}
