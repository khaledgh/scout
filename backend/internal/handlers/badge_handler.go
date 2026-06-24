package handlers

import (
	"kashfi/internal/dto"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BadgeHandler struct {
	svc    *services.BadgeService
	gamify *services.GamificationService
}

func NewBadgeHandler(svc *services.BadgeService, gamify *services.GamificationService) *BadgeHandler {
	return &BadgeHandler{svc: svc, gamify: gamify}
}

func (h *BadgeHandler) ListBadges(c echo.Context) error {
	badges, err := h.svc.ListCatalog(c.Request().Context())
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, badges)
}

func (h *BadgeHandler) CreateBadge(c echo.Context) error {
	var req dto.CreateBadgeRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	badge, err := h.svc.Create(c.Request().Context(), services.CreateBadgeInput{
		Name: req.Name, Description: req.Description, Category: req.Category, XPReward: req.XPReward,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, badge)
}

func (h *BadgeHandler) GetBadge(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	badge, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return utils.NotFound(c, "badge")
	}
	return utils.OK(c, badge)
}

func (h *BadgeHandler) UpdateBadge(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.UpdateBadgeRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	badge, err := h.svc.Update(c.Request().Context(), id, services.UpdateBadgeInput{
		Name: req.Name, Description: req.Description, Category: req.Category,
		XPReward: req.XPReward, IsActive: req.IsActive,
	})
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "badge") }
		return utils.Internal(c)
	}
	return utils.OK(c, badge)
}

func (h *BadgeHandler) DeleteBadge(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "badge") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *BadgeHandler) RevokeBadge(c echo.Context) error {
	memberID, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid member id")
	}
	badgeID, err := strconv.ParseUint(c.Param("badgeId"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "invalid badge id")
	}
	if err := h.svc.RevokeBadge(c.Request().Context(), memberID, uint(badgeID)); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "member badge") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *BadgeHandler) AwardBadge(c echo.Context) error {
	memberID, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid member id")
	}
	var req dto.AwardBadgeRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	awardedBy := appMiddleware.GetUserID(c)
	mb, err := h.svc.AwardBadge(c.Request().Context(), memberID, req.BadgeID, awardedBy)
	if err != nil {
		if err == services.ErrConflict { return utils.Conflict(c, "badge already awarded") }
		if err == services.ErrNotFound { return utils.NotFound(c, "badge") }
		return utils.Internal(c)
	}
	return utils.Created(c, mb)
}

func (h *BadgeHandler) MemberBadges(c echo.Context) error {
	memberID, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid member id")
	}
	badges, err := h.svc.MemberBadges(c.Request().Context(), memberID)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, badges)
}

func (h *BadgeHandler) ListSkills(c echo.Context) error {
	skills, err := h.svc.ListSkills(c.Request().Context())
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, skills)
}

func (h *BadgeHandler) CreateSkill(c echo.Context) error {
	var req dto.CreateSkillRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	maxLevel := req.MaxLevel
	if maxLevel == 0 { maxLevel = 5 }
	sk, err := h.svc.CreateSkill(c.Request().Context(), req.Name, req.Category, req.Description, maxLevel)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, sk)
}

func (h *BadgeHandler) UpdateSkill(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.UpdateSkillRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	sk, err := h.svc.UpdateSkill(c.Request().Context(), id, services.UpdateSkillInput{
		Name: req.Name, Category: req.Category, Description: req.Description, MaxLevel: req.MaxLevel,
	})
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "skill") }
		return utils.Internal(c)
	}
	return utils.OK(c, sk)
}

func (h *BadgeHandler) DeleteSkill(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if err := h.svc.DeleteSkill(c.Request().Context(), id); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "skill") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *BadgeHandler) AssessSkill(c echo.Context) error {
	memberID, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid member id")
	}
	var req dto.AssessSkillRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	assessorID := appMiddleware.GetUserID(c)
	ms, err := h.svc.AssessSkill(c.Request().Context(), memberID, req.SkillID, assessorID, req.Level)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, ms)
}

func (h *BadgeHandler) MemberLeaderboard(c echo.Context) error {
	section := c.QueryParam("section")
	var sectionPtr *string
	if section != "" { sectionPtr = &section }
	result, err := h.gamify.MemberLeaderboard(c.Request().Context(), sectionPtr, 50)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, result)
}

func (h *BadgeHandler) UnitLeaderboard(c echo.Context) error {
	result, err := h.gamify.UnitLeaderboard(c.Request().Context())
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, result)
}

func (h *BadgeHandler) XPHistory(c echo.Context) error {
	memberID, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid member id")
	}
	events, err := h.gamify.GetXPHistory(c.Request().Context(), memberID)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, events)
}
