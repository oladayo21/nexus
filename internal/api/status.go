package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const appVersion = "0.1.0"

type StatusResponse struct {
	SetupRequired bool          `json:"setupRequired"`
	Authenticated bool          `json:"authenticated"`
	User          *UserResponse `json:"user"`
	Version       string        `json:"version"`
	Env           string        `json:"env"`
}

func (s *APIServer) handleStatus(c echo.Context) error {
	resp := StatusResponse{
		Version: appVersion,
		Env:     s.opts.Config.Env,
	}

	hasUsers, err := s.db.Queries().HasUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check setup status")
	}

	resp.SetupRequired = !hasUsers

	userID, ok := s.GetUserID(c)
	if ok {
		pgUUID, err := stringToUUID(userID)
		if err == nil {
			user, err := s.db.Queries().GetUserByID(c.Request().Context(), pgUUID)
			if err == nil {
				resp.Authenticated = true
				resp.User = &UserResponse{
					ID:    uuidToString(user.ID),
					Name:  user.Name,
					Email: user.Email,
				}
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}
