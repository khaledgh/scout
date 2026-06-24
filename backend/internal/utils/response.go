package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

func OKWithMeta(c echo.Context, data interface{}, meta *Meta) error {
	return c.JSON(http.StatusOK, Response{Success: true, Data: data, Meta: meta})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{Success: true, Data: data})
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func Fail(c echo.Context, status int, code, message string) error {
	return c.JSON(status, Response{
		Success: false,
		Error:   &APIError{Code: code, Message: message},
	})
}

func BadRequest(c echo.Context, message string) error {
	return Fail(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(c echo.Context) error {
	return Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
}

func Forbidden(c echo.Context) error {
	return Fail(c, http.StatusForbidden, "FORBIDDEN", "insufficient permissions")
}

func NotFound(c echo.Context, resource string) error {
	return Fail(c, http.StatusNotFound, "NOT_FOUND", resource+" not found")
}

func Conflict(c echo.Context, message string) error {
	return Fail(c, http.StatusConflict, "CONFLICT", message)
}

func Internal(c echo.Context) error {
	return Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
}
