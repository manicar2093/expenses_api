package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/middlewares"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/recurrent"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
)

type RecurrentExpensesController struct {
	middlewares.Middlewares
	recurrentExpenseCreator recurrent.RecurrentExpenseCreator
	group                   *echo.Group
}

func NewRecurrentExpensesController(
	recurrentExpenseCreator recurrent.RecurrentExpenseCreator,
	middlewares middlewares.Middlewares,
	e *echo.Echo,
) *RecurrentExpensesController {
	c := &RecurrentExpensesController{recurrentExpenseCreator: recurrentExpenseCreator, Middlewares: middlewares, group: e.Group("/recurrent_expenses")}
	c.register()
	return c
}

func (c *RecurrentExpensesController) register() {
	c.group.POST("", c.Create, c.LoggedIn)
}

// @Summary     Create a recurrent expense
// @Description Register a recurrent expense into the database
// @Tags        recurrent expenses
// @Accept      json
// @Produce     json
// @Param       expense_to_create body     entities.RecurrentExpense true "Expense to be created"
// @Success     201               {object} entities.RecurrentExpense "Expense has been created"
// @Failure     400               {object} validator.ValidationError "When a request does not fulfill need data"
// @Failure     500               "Something unidentified has occurred"
// @Security    ApiKeyAuth
// @Router      /recurrent_expenses [post]
func (c *RecurrentExpensesController) Create(ctx echo.Context) error {
	var request entities.RecurrentExpense
	if err := ctx.Bind(&request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	request.UserID = getUserIDAsUUID(ctx)
	if err := ctx.Validate(request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	if err := c.recurrentExpenseCreator.Create(ctx.Request().Context(), &request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, &request)
}
