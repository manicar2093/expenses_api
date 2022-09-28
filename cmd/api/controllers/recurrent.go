package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/pkg/errors"
)

type RecurrentExpensesController struct {
	createRecurrentExpense         recurrentexpenses.CreateRecurrentExpense
	getAllRecurrentExpenses        recurrentexpenses.GetAllRecurrentExpenses
	createMonthlyRecurrentExpenses recurrentexpenses.CreateMonthlyRecurrentExpenses
	group                          *echo.Group
}

func NewRecurrentExpensesController(
	createRecurrentExpense recurrentexpenses.CreateRecurrentExpense,
	getAllRecurrentExpenses recurrentexpenses.GetAllRecurrentExpenses,
	createMonthlyRecurrentExpenses recurrentexpenses.CreateMonthlyRecurrentExpenses,
	e *echo.Echo, //nolint:varnamelen
) *RecurrentExpensesController {
	return &RecurrentExpensesController{
		createRecurrentExpense:         createRecurrentExpense,
		getAllRecurrentExpenses:        getAllRecurrentExpenses,
		createMonthlyRecurrentExpenses: createMonthlyRecurrentExpenses,
		group:                          e.Group("/recurrent_expenses"),
	}
}

func (c *RecurrentExpensesController) Register() {
	c.group.POST("", c.create)

	c.group.GET("/all", c.getAll)

	c.group.POST("/monthly_expenses", c.createMonthly)
}

func (c *RecurrentExpensesController) create(ctx echo.Context) error {
	var recurrentExpenseReq recurrentexpenses.CreateRecurrentExpenseInput
	if err := ctx.Bind(&recurrentExpenseReq); err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	res, err := c.createRecurrentExpense.Create(ctx.Request().Context(), &recurrentExpenseReq)
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, res)
}

func (c *RecurrentExpensesController) getAll(ctx echo.Context) error {
	res, err := c.getAllRecurrentExpenses.GetAll(ctx.Request().Context())
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, res)
}

func (c *RecurrentExpensesController) createMonthly(ctx echo.Context) error {
	got, err := c.createMonthlyRecurrentExpenses.CreateMonthlyRecurrentExpenses(ctx.Request().Context())
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, got)
}
