package handlers

import (
	"kashfi/internal/dto"
	"kashfi/internal/models"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UnitHandler struct {
	svc *services.UnitService
}

func NewUnitHandler(svc *services.UnitService) *UnitHandler {
	return &UnitHandler{svc: svc}
}

func (h *UnitHandler) List(c echo.Context) error {
	units, err := h.svc.List(c.Request().Context())
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, units)
}

func (h *UnitHandler) Get(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	unit, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return utils.NotFound(c, "unit")
	}
	return utils.OK(c, unit)
}

func (h *UnitHandler) Create(c echo.Context) error {
	var req dto.CreateUnitRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	unit, err := h.svc.Create(c.Request().Context(), services.CreateUnitInput{
		Name: req.Name, Section: models.Section(req.Section), Motto: req.Motto,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, unit)
}

func (h *UnitHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.UpdateUnitRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	in := services.UpdateUnitInput{}
	if req.Name != nil { in.Name = req.Name }
	if req.Motto != nil { in.Motto = req.Motto }
	if req.IsActive != nil { in.IsActive = req.IsActive }
	if req.Section != nil { s := models.Section(*req.Section); in.Section = &s }

	unit, err := h.svc.Update(c.Request().Context(), id, in)
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "unit") }
		return utils.Internal(c)
	}
	return utils.OK(c, unit)
}

func (h *UnitHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "unit") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *UnitHandler) AddMembers(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.AddUnitMembersRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := h.svc.AddMembers(c.Request().Context(), id, req.MemberIDs); err != nil {
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *UnitHandler) RemoveMember(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	mid, err := strconv.ParseUint(c.Param("mid"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "invalid member id")
	}
	if err := h.svc.RemoveMember(c.Request().Context(), id, uint(mid)); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "unit member") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *UnitHandler) AssignLeader(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.AssignUnitLeaderRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	if err := h.svc.AssignLeader(c.Request().Context(), id, req.UserID, models.UnitLeaderRole(req.RoleInUnit)); err != nil {
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *UnitHandler) Leaderboard(c echo.Context) error {
	units, err := h.svc.Leaderboard(c.Request().Context())
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, units)
}
