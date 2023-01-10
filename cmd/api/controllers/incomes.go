package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/middlewares"
	"github.com/manicar2093/expenses_api/internal/incomes"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/goption"
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
	var incomeData incomes.CreateIncomeInput
	if err := ctx.Bind(&incomeData); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	incomeData.UserID = getUserIDAsOptionlaUUID(ctx)
	got, err := c.CreateIncome.Create(ctx.Request().Context(), &incomeData)
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, &got)
}

func getUserIDAsOptionlaUUID(ctx echo.Context) goption.Optional[uuid.UUID] {
	return goption.Of(uuid.MustParse(ctx.Get("user_id").(string)))
}
