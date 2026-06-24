package middleware

import (
	"kashfi/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "user_role"
)

func JWT(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				return utils.Unauthorized(c)
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := utils.ParseToken(tokenStr, secret)
			if err != nil {
				return utils.Unauthorized(c)
			}

			if claims.Type != "access" {
				return utils.Unauthorized(c)
			}

			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyRole, claims.Role)
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) uint {
	if v, ok := c.Get(ContextKeyUserID).(uint); ok {
		return v
	}
	return 0
}

func GetRole(c echo.Context) string {
	if v, ok := c.Get(ContextKeyRole).(string); ok {
		return v
	}
	return ""
}
