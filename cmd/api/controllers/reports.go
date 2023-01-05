package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/middlewares"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
)

type ReportsController struct {
	middlewares.Middlewares
	getCurrentMonth reports.CurrentMonthDetailsGettable
	group           *echo.Group
}

func NewReportsController(getCurrentMonth reports.CurrentMonthDetailsGettable, middlewares middlewares.Middlewares, e *echo.Echo) *ReportsController {
	controller := &ReportsController{
		Middlewares:     middlewares,
		getCurrentMonth: getCurrentMonth,
		group:           e.Group("/reports"),
	}
	controller.register()
	return controller
}

func (c *ReportsController) register() {
	c.group.GET("/current_month", c.CurrentMonth, c.LoggedIn)
}

// @Summary     Get current month details
// @Description Generates current month general details
// @Tags        reports
// @Produce     json
// @Success     200
// @Failure     500
// @Security    ApiKeyAuth
// @Router      /reports/current_month [get]
func (c *ReportsController) CurrentMonth(ctx echo.Context) error {
	currentMonthDetails, err := c.getCurrentMonth.GetExpenses(ctx.Request().Context(), uuid.MustParse(ctx.Get("user_id").(string)))
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, currentMonthDetails)
}
