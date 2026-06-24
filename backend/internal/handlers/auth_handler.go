package handlers

import (
	"kashfi/internal/dto"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	result, err := h.svc.Login(c.Request().Context(), req.Phone, req.Password)
	if err != nil {
		return utils.Fail(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid phone or password")
	}
	return utils.OK(c, echo.Map{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"user":          result.User,
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req dto.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	result, err := h.svc.Refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return utils.Unauthorized(c)
	}
	return utils.OK(c, echo.Map{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	return utils.NoContent(c)
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID := appMiddleware.GetUserID(c)
	user, err := h.svc.Me(c.Request().Context(), userID)
	if err != nil {
		return utils.NotFound(c, "user")
	}
	return utils.OK(c, user)
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	userID := appMiddleware.GetUserID(c)
	if err := h.svc.ChangePassword(c.Request().Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		if err == services.ErrUnauthorized {
			return utils.Fail(c, http.StatusUnauthorized, "WRONG_PASSWORD", "current password is incorrect")
		}
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}
