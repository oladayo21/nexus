package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/oladayo21/nexus/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// Session helpers

func (s *APIServer) GetUserID(c echo.Context) (string, bool) {
	userID := s.sessionManager.GetString(c.Request().Context(), "userID")

	return userID, userID != ""
}

func (s *APIServer) SetUserID(c echo.Context, userID string) error {
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
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		return next(c)
	}
}

// Password helpers

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Request/Response types

type SignupRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UUID conversion helpers

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}

	return uuid.UUID(u.Bytes).String()
}

func stringToUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

// Auth handlers

func (s *APIServer) handleSignup(c echo.Context) error {
	var req SignupRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if err := c.Validate(&req); err != nil {
		return err
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	user, err := s.db.Queries().CreateUser(c.Request().Context(), database.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return echo.NewHTTPError(http.StatusConflict, "email already exists")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	userID := uuidToString(user.ID)

	if err := s.SetUserID(c, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}

	return c.JSON(http.StatusCreated, UserResponse{
		ID:    userID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (s *APIServer) handleLogin(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if err := c.Validate(&req); err != nil {
		return err
	}

	user, err := s.db.Queries().GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	if !checkPassword(user.PasswordHash, req.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	userID := uuidToString(user.ID)

	if err := s.SetUserID(c, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}

	return c.JSON(http.StatusOK, UserResponse{
		ID:    userID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (s *APIServer) handleLogout(c echo.Context) error {
	if err := s.ClearSession(c); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to clear session")
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *APIServer) handleMe(c echo.Context) error {
	userID, _ := s.GetUserID(c)

	pgUUID, err := stringToUUID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid session")
	}

	user, err := s.db.Queries().GetUserByID(c.Request().Context(), pgUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, UserResponse{
		ID:    uuidToString(user.ID),
		Name:  user.Name,
		Email: user.Email,
	})
}
