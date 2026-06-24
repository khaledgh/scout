package handlers

import (
	"kashfi/internal/dto"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/models"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type MemberHandler struct {
	svc    *services.MemberService
	qrSecret string
}

func NewMemberHandler(svc *services.MemberService, qrSecret string) *MemberHandler {
	return &MemberHandler{svc: svc, qrSecret: qrSecret}
}

func (h *MemberHandler) List(c echo.Context) error {
	p := utils.ParsePagination(c)
	f := services.MemberFilter{}
	if v := c.QueryParam("unit_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			uid := uint(id); f.UnitID = &uid
		}
	}
	if v := c.QueryParam("section"); v != "" { f.Section = &v }
	if v := c.QueryParam("status"); v != "" { f.Status = &v }
	if v := c.QueryParam("search"); v != "" { f.Search = &v }

	members, total, err := h.svc.List(c.Request().Context(), f, p.Page, p.PageSize)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OKWithMeta(c, members, utils.BuildMeta(p, total))
}

func (h *MemberHandler) Get(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	m, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return utils.NotFound(c, "member")
	}
	return utils.OK(c, m)
}

func (h *MemberHandler) Create(c echo.Context) error {
	var req dto.CreateMemberRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	bd, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return utils.BadRequest(c, "invalid birth_date format, use YYYY-MM-DD")
	}
	jd, err := time.Parse("2006-01-02", req.JoinDate)
	if err != nil {
		return utils.BadRequest(c, "invalid join_date format, use YYYY-MM-DD")
	}
	m, err := h.svc.Create(c.Request().Context(), services.CreateMemberInput{
		FullName: req.FullName, BirthDate: bd, Gender: req.Gender,
		Section: models.Section(req.Section), RankStage: req.RankStage, JoinDate: jd,
		ParentName: req.ParentName, ParentPhone: req.ParentPhone,
		SecondPhone: req.SecondPhone, Address: req.Address, UserID: req.UserID,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, m)
}

func (h *MemberHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if !h.svc.CanManageMember(c.Request().Context(), appMiddleware.GetUserID(c), appMiddleware.GetRole(c), id) {
		return utils.Forbidden(c)
	}
	var req dto.UpdateMemberRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	in := services.UpdateMemberInput{}
	if req.FullName != nil { in.FullName = req.FullName }
	if req.Gender != nil { in.Gender = req.Gender }
	if req.RankStage != nil { in.RankStage = req.RankStage }
	if req.ParentName != nil { in.ParentName = req.ParentName }
	if req.ParentPhone != nil { in.ParentPhone = req.ParentPhone }
	if req.SecondPhone != nil { in.SecondPhone = req.SecondPhone }
	if req.Address != nil { in.Address = req.Address }
	if req.Section != nil { s := models.Section(*req.Section); in.Section = &s }
	if req.Status != nil { st := models.MemberStatus(*req.Status); in.Status = &st }
	if req.BirthDate != nil {
		if bd, err := time.Parse("2006-01-02", *req.BirthDate); err == nil { in.BirthDate = &bd }
	}
	if req.JoinDate != nil {
		if jd, err := time.Parse("2006-01-02", *req.JoinDate); err == nil { in.JoinDate = &jd }
	}

	m, err := h.svc.Update(c.Request().Context(), id, in)
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "member") }
		return utils.Internal(c)
	}
	return utils.OK(c, m)
}

func (h *MemberHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if !h.svc.CanManageMember(c.Request().Context(), appMiddleware.GetUserID(c), appMiddleware.GetRole(c), id) {
		return utils.Forbidden(c)
	}
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "member") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *MemberHandler) GetMedical(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if !h.svc.CanManageMember(c.Request().Context(), appMiddleware.GetUserID(c), appMiddleware.GetRole(c), id) {
		return utils.Forbidden(c)
	}
	med, err := h.svc.GetMedical(c.Request().Context(), id)
	if err != nil {
		return utils.NotFound(c, "medical record")
	}
	return utils.OK(c, med)
}

func (h *MemberHandler) UpsertMedical(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if !h.svc.CanManageMember(c.Request().Context(), appMiddleware.GetUserID(c), appMiddleware.GetRole(c), id) {
		return utils.Forbidden(c)
	}
	var req dto.UpsertMedicalRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	med, err := h.svc.UpsertMedical(c.Request().Context(), id, services.MedicalInput{
		BloodType: req.BloodType, Allergies: req.Allergies,
		ChronicConditions: req.ChronicConditions, Medications: req.Medications,
		EmergencyNotes: req.EmergencyNotes,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, med)
}

func (h *MemberHandler) Timeline(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	tl, err := h.svc.GetTimeline(c.Request().Context(), id)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, tl)
}

func (h *MemberHandler) CreateEvaluation(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if !h.svc.CanManageMember(c.Request().Context(), appMiddleware.GetUserID(c), appMiddleware.GetRole(c), id) {
		return utils.Forbidden(c)
	}
	var req dto.CreateEvaluationRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	evaluatorID := appMiddleware.GetUserID(c)
	eval, err := h.svc.CreateEvaluation(c.Request().Context(), id, evaluatorID, services.EvalInput{
		Period: req.Period, Discipline: req.Discipline, Participation: req.Participation,
		Leadership: req.Leadership, Skill: req.Skill, Overall: req.Overall, Notes: req.Notes,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, eval)
}

func (h *MemberHandler) UploadPhoto(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if !h.svc.CanManageMember(c.Request().Context(), appMiddleware.GetUserID(c), appMiddleware.GetRole(c), id) {
		return utils.Forbidden(c)
	}
	file, err := c.FormFile("file")
	if err != nil {
		return utils.BadRequest(c, "file required")
	}
	filename, err := utils.SaveUpload(file, h.svc.UploadPath(), h.svc.UploadTypes(), h.svc.UploadMaxMB())
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	m, err := h.svc.SetPhoto(c.Request().Context(), id, filename)
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "member") }
		return utils.Internal(c)
	}
	return utils.OK(c, m)
}

func (h *MemberHandler) GetQRToken(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	token, err := h.svc.GetQRToken(c.Request().Context(), id, h.qrSecret)
	if err != nil {
		return utils.NotFound(c, "member")
	}
	return utils.OK(c, echo.Map{"token": token, "member_id": id})
}

func parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err
}
