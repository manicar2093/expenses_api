package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/middlewares"
	"github.com/manicar2093/expenses_api/internal/incomes"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
)

type IncomesController struct {
	middlewares.Middlewares
	incomes.CreateIncome
	*echo.Group
}

func NewIncomesController(middlewares middlewares.Middlewares, incomesCreator incomes.CreateIncome, e *echo.Echo) *IncomesController {
	c := &IncomesController{ //nolint:varnamelen
		Middlewares:  middlewares,
		CreateIncome: incomesCreator,
		Group:        e.Group("/incomes"),
	}
	c.register()
	return c
}

func (c *IncomesController) register() {
	c.POST("", c.Create, c.LoggedIn)
}

// @Summary     Create an income
// @Description Register a income into the database
// @Tags        incomes
// @Accept      json
// @Produce     json
// @Param       expense_to_create body     incomes.CreateIncomeInput true "Income to be created"
// @Success     201               {object} entities.Income           "Income has been created"
// @Failure     400               {object} validator.ValidationError "When a request does not fulfill need data"
// @Failure     500               "Something unidentified has occurred"
// @Security    ApiKeyAuth
// @Router      /incomes [post]
func (c *IncomesController) Create(ctx echo.Context) error {
	var request incomes.CreateIncomeInput
	if err := ctx.Bind(&request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	request.UserID = getUserIDAsOptionlaUUID(ctx)
	if err := ctx.Validate(request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	got, err := c.CreateIncome.Create(ctx.Request().Context(), &request)
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, &got)
}
