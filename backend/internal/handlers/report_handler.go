package handlers

import (
	"encoding/csv"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	svc *services.ReportService
}

func NewReportHandler(svc *services.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) Dashboard(c echo.Context) error {
	data, err := h.svc.Dashboard(c.Request().Context())
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, data)
}

func (h *ReportHandler) MemberReport(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	report, err := h.svc.MemberReport(c.Request().Context(), id)
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "member") }
		return utils.Internal(c)
	}
	return utils.OK(c, report)
}

func (h *ReportHandler) UnitReport(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	report, err := h.svc.UnitReport(c.Request().Context(), id)
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "unit") }
		return utils.Internal(c)
	}
	return utils.OK(c, report)
}

func (h *ReportHandler) Monthly(c echo.Context) error {
	monthStr := c.QueryParam("month") // format: 2026-06
	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		t = time.Now()
	}
	report, err := h.svc.Monthly(c.Request().Context(), t.Year(), int(t.Month()))
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, report)
}

func (h *ReportHandler) Export(c echo.Context) error {
	reportType := c.QueryParam("type")
	month := c.QueryParam("month")
	t, _ := time.Parse("2006-01", month)
	if t.IsZero() { t = time.Now() }

	switch reportType {
	case "monthly":
		report, err := h.svc.Monthly(c.Request().Context(), t.Year(), int(t.Month()))
		if err != nil {
			return utils.Internal(c)
		}
		c.Response().Header().Set("Content-Disposition", `attachment; filename="monthly-report.csv"`)
		c.Response().Header().Set("Content-Type", "text/csv")
		c.Response().WriteHeader(http.StatusOK)
		w := csv.NewWriter(c.Response())
		w.Write([]string{"Month", "Activities", "Attendance Rate %", "New Members", "XP Distributed"})
		w.Write([]string{
			report.Month,
			strconv.FormatInt(report.TotalActivities, 10),
			strconv.FormatFloat(report.AttendanceRate, 'f', 1, 64),
			strconv.FormatInt(report.NewMembers, 10),
			strconv.FormatInt(report.XPDistributed, 10),
		})
		w.Flush()
		return nil
	default:
		return utils.BadRequest(c, "unsupported report type")
	}
}
