package handlers

import (
	"kashfi/internal/dto"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/models"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"mime/multipart"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type ActivityHandler struct {
	svc        *services.ActivityService
	attendSvc  *services.AttendanceService
}

func NewActivityHandler(svc *services.ActivityService, attendSvc *services.AttendanceService) *ActivityHandler {
	return &ActivityHandler{svc: svc, attendSvc: attendSvc}
}

func (h *ActivityHandler) List(c echo.Context) error {
	p := utils.ParsePagination(c)
	f := services.ActivityFilter{}
	if v := c.QueryParam("type"); v != "" { f.Type = &v }
	if v := c.QueryParam("status"); v != "" { f.Status = &v }
	if v := c.QueryParam("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil { f.From = &t }
	}
	if v := c.QueryParam("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil { f.To = &t }
	}
	activities, total, err := h.svc.List(c.Request().Context(), f, p.Page, p.PageSize)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OKWithMeta(c, activities, utils.BuildMeta(p, total))
}

func (h *ActivityHandler) Get(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	a, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return utils.NotFound(c, "activity")
	}
	return utils.OK(c, a)
}

func (h *ActivityHandler) Create(c echo.Context) error {
	var req dto.CreateActivityRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	starts, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		starts, err = time.Parse("2006-01-02T15:04", req.StartsAt)
		if err != nil {
			return utils.BadRequest(c, "invalid starts_at")
		}
	}
	ends, err := time.Parse(time.RFC3339, req.EndsAt)
	if err != nil {
		ends, err = time.Parse("2006-01-02T15:04", req.EndsAt)
		if err != nil {
			return utils.BadRequest(c, "invalid ends_at")
		}
	}
	a, err := h.svc.Create(c.Request().Context(), services.CreateActivityInput{
		Title: req.Title, Description: req.Description, Type: models.ActivityType(req.Type),
		Location: req.Location, LocationLat: req.LocationLat, LocationLng: req.LocationLng,
		StartsAt: starts, EndsAt: ends,
		ResponsibleUserID: appMiddleware.GetUserID(c), UnitID: req.UnitID,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, a)
}

func (h *ActivityHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.UpdateActivityRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	in := services.UpdateActivityInput{}
	if req.Title != nil { in.Title = req.Title }
	if req.Description != nil { in.Description = req.Description }
	if req.Location != nil { in.Location = req.Location }
	if req.LocationLat != nil { in.LocationLat = req.LocationLat }
	if req.LocationLng != nil { in.LocationLng = req.LocationLng }
	if req.UnitID != nil { in.UnitID = req.UnitID }
	if req.Type != nil { t := models.ActivityType(*req.Type); in.Type = &t }
	if req.Status != nil { s := models.ActivityStatus(*req.Status); in.Status = &s }

	a, err := h.svc.Update(c.Request().Context(), id, in)
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "activity") }
		return utils.Internal(c)
	}
	return utils.OK(c, a)
}

func (h *ActivityHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "activity") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *ActivityHandler) UploadMedia(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	file, err := c.FormFile("file")
	if err != nil {
		return utils.BadRequest(c, "file required")
	}
	filename, err := utils.SaveUpload(file, h.svc.UploadPath(), h.svc.UploadTypes(), h.svc.UploadMaxMB())
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	mediaType := detectMediaType(file)
	media, err := h.svc.AddMedia(c.Request().Context(), id, appMiddleware.GetUserID(c), filename, mediaType)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, media)
}

func (h *ActivityHandler) GetAttendance(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	records, err := h.attendSvc.GetForActivity(c.Request().Context(), id)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, records)
}

func (h *ActivityHandler) RecordAttendance(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.RecordAttendanceRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	records := make([]services.AttendanceRecord, len(req.Records))
	for i, r := range req.Records {
		records[i] = services.AttendanceRecord{
			MemberID: r.MemberID,
			Status:   models.AttendanceStatus(r.Status),
		}
	}
	if err := h.attendSvc.BulkRecord(c.Request().Context(), id, appMiddleware.GetUserID(c), records); err != nil {
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *ActivityHandler) CheckIn(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.CheckInRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	userID := appMiddleware.GetUserID(c)

	switch req.Method {
	case "qr":
		if req.QRToken == nil {
			return utils.BadRequest(c, "qr_token required")
		}
		att, err := h.attendSvc.CheckInQR(c.Request().Context(), id, *req.QRToken, userID)
		if err != nil {
			if err == services.ErrBadRequest { return utils.BadRequest(c, "invalid or expired QR token") }
			return utils.Internal(c)
		}
		return utils.OK(c, att)
	case "gps":
		if req.Lat == nil || req.Lng == nil {
			return utils.BadRequest(c, "lat and lng required for GPS check-in")
		}
		att, err := h.attendSvc.CheckInGPS(c.Request().Context(), id, userID, *req.Lat, *req.Lng, userID)
		if err != nil {
			if err == services.ErrForbidden { return utils.Fail(c, 403, "OUT_OF_RANGE", "you are outside the geofence radius") }
			if err == services.ErrBadRequest { return utils.BadRequest(c, "activity has no location set") }
			return utils.Internal(c)
		}
		return utils.OK(c, att)
	default:
		return utils.BadRequest(c, "invalid check-in method")
	}
}

func (h *ActivityHandler) CreateFeedback(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.FeedbackRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	memberID := appMiddleware.GetUserID(c)
	fb, err := h.svc.CreateFeedback(c.Request().Context(), id, memberID, req.Rating, req.WhatWentWell, req.WhatToImprove, req.Comment)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, fb)
}

func (h *ActivityHandler) FeedbackSummary(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	summary, err := h.svc.FeedbackSummary(c.Request().Context(), id)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, summary)
}

func detectMediaType(file *multipart.FileHeader) models.MediaType {
	ct := file.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "video/") {
		return models.MediaVideo
	}
	return models.MediaImage
}
