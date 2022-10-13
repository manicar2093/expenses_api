package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/pkg/errors"
)

type ReportsController struct {
	getCurrentMonth reports.CurrentMonthDetailsGettable
	group           *echo.Group
}

func NewReportsController(getCurrentMonth reports.CurrentMonthDetailsGettable, e *echo.Echo) *ReportsController {
	return &ReportsController{getCurrentMonth: getCurrentMonth, group: e.Group("/reports")}
}

func (c *ReportsController) Register() {
	c.group.GET("/current_month", c.currentMonth)
}

// @Summary     Get current month details
// @Description Generates current month general details
// @Tags        reports
// @Produce     json
// @Success     200
// @Failure     500
// @Router      /reports/current_month [get]
func (c *ReportsController) currentMonth(ctx echo.Context) error {
	currentMonthDetails, err := c.getCurrentMonth.GetExpenses(ctx.Request().Context())
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, currentMonthDetails)
}
