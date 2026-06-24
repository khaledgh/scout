package utils

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Pagination struct {
	Page     int
	PageSize int
	Offset   int
}

func ParsePagination(c echo.Context) Pagination {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

func BuildMeta(p Pagination, total int64) *Meta {
	totalPages := int(math.Ceil(float64(total) / float64(p.PageSize)))
	return &Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
