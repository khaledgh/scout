package handlers

import (
	"kashfi/internal/dto"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/services"
	"kashfi/internal/utils"

	"github.com/labstack/echo/v4"
)

type AnnouncementHandler struct {
	svc   *services.AnnouncementService
	notif *services.NotificationService
}

func NewAnnouncementHandler(svc *services.AnnouncementService, notif *services.NotificationService) *AnnouncementHandler {
	return &AnnouncementHandler{svc: svc, notif: notif}
}

func (h *AnnouncementHandler) List(c echo.Context) error {
	userID := appMiddleware.GetUserID(c)
	role := appMiddleware.GetRole(c)
	announcements, err := h.svc.List(c.Request().Context(), userID, role)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, announcements)
}

func (h *AnnouncementHandler) Create(c echo.Context) error {
	var req dto.CreateAnnouncementRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	authorID := appMiddleware.GetUserID(c)
	a, err := h.svc.Create(c.Request().Context(), services.CreateAnnouncementInput{
		Title: req.Title, Body: req.Body, Audience: req.Audience,
		UnitID: req.UnitID, Pinned: req.Pinned, AuthorID: authorID,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, a)
}

func (h *AnnouncementHandler) Get(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	a, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return utils.NotFound(c, "announcement")
	}
	return utils.OK(c, a)
}

func (h *AnnouncementHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.UpdateAnnouncementRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	a, err := h.svc.Update(c.Request().Context(), id, appMiddleware.GetUserID(c), appMiddleware.GetRole(c), services.UpdateAnnouncementInput{
		Title: req.Title, Body: req.Body, Audience: req.Audience, UnitID: req.UnitID, Pinned: req.Pinned,
	})
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "announcement") }
		if err == services.ErrForbidden { return utils.Forbidden(c) }
		return utils.Internal(c)
	}
	return utils.OK(c, a)
}

func (h *AnnouncementHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if err := h.svc.Delete(c.Request().Context(), id, appMiddleware.GetUserID(c), appMiddleware.GetRole(c)); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "announcement") }
		if err == services.ErrForbidden { return utils.Forbidden(c) }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

type NotificationHandler struct {
	svc *services.NotificationService
}

func NewNotificationHandler(svc *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) List(c echo.Context) error {
	userID := appMiddleware.GetUserID(c)
	notifs, err := h.svc.ForUser(c.Request().Context(), userID)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, notifs)
}

func (h *NotificationHandler) MarkRead(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	userID := appMiddleware.GetUserID(c)
	if err := h.svc.MarkRead(c.Request().Context(), id, userID); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "notification") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}
