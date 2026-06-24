package middleware

import (
	"kashfi/internal/utils"

	"github.com/labstack/echo/v4"
)

func RequireRole(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := GetRole(c)
			if !allowed[role] {
				return utils.Forbidden(c)
			}
			return next(c)
		}
	}
}

// IsAdmin returns true for super_admin role.
func IsAdmin(c echo.Context) bool {
	return GetRole(c) == "super_admin"
}

// IsLeaderOrAdmin returns true for leader/super_admin.
func IsLeaderOrAdmin(c echo.Context) bool {
	r := GetRole(c)
	return r == "leader" || r == "super_admin" || r == "assistant"
}
