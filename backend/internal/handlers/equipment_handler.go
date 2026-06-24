package handlers

import (
	"kashfi/internal/dto"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type EquipmentHandler struct {
	svc *services.EquipmentService
}

func NewEquipmentHandler(svc *services.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{svc: svc}
}

func (h *EquipmentHandler) List(c echo.Context) error {
	items, err := h.svc.List(c.Request().Context())
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, items)
}

func (h *EquipmentHandler) Get(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	item, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return utils.NotFound(c, "equipment")
	}
	return utils.OK(c, item)
}

func (h *EquipmentHandler) Create(c echo.Context) error {
	var req dto.CreateEquipmentRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	item, err := h.svc.Create(c.Request().Context(), services.CreateEquipmentInput{
		Name: req.Name, Category: req.Category, QuantityTotal: req.QuantityTotal,
		QuantityAvailable: req.QuantityAvailable, Condition: req.Condition, Notes: req.Notes,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, item)
}

func (h *EquipmentHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.UpdateEquipmentRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	item, err := h.svc.Update(c.Request().Context(), id, services.UpdateEquipmentInput{
		Name: req.Name, Category: req.Category, QuantityTotal: req.QuantityTotal,
		QuantityAvailable: req.QuantityAvailable, Condition: req.Condition, Notes: req.Notes,
	})
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "equipment") }
		return utils.Internal(c)
	}
	return utils.OK(c, item)
}

func (h *EquipmentHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "equipment") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *EquipmentHandler) Loans(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	loans, err := h.svc.Loans(c.Request().Context(), id)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, loans)
}

func (h *EquipmentHandler) Loan(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.LoanEquipmentRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	due, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return utils.BadRequest(c, "invalid due_date, use YYYY-MM-DD")
	}
	loan, err := h.svc.Loan(c.Request().Context(), services.LoanInput{
		EquipmentID: id, BorrowedBy: req.BorrowedBy, ActivityID: req.ActivityID,
		Quantity: req.Quantity, DueDate: due,
	})
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "equipment") }
		if err == services.ErrBadRequest { return utils.BadRequest(c, "insufficient quantity available") }
		return utils.Internal(c)
	}
	return utils.Created(c, loan)
}

func (h *EquipmentHandler) ReturnLoan(c echo.Context) error {
	loanID, err := strconv.ParseUint(c.Param("loanId"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "invalid loan id")
	}
	if err := h.svc.ReturnLoan(c.Request().Context(), uint(loanID)); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "loan") }
		if err == services.ErrBadRequest { return utils.BadRequest(c, "loan already returned") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}
