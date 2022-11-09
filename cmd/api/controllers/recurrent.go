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

// @Summary     Create a recurrent expense
// @Description Register into db a recurrent expense
// @Tags        recurrent expenses
// @Accept      json
// @Produce     json
// @Param       recurrent_expense body recurrentexpenses.CreateRecurrentExpenseInput true "Recurrent Expense"
// @Success     201 {object} recurrentexpenses.CreateRecurrentExpenseOutput
// @Failure     500
// @Router      /recurrent_expenses [post]
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

// @Summary     Get all recurrent expense
// @Description Finds all recurrent expenses registered in db
// @Tags        recurrent expenses
// @Produce     json
// @Success     200 {object} recurrentexpenses.GetAllRecurrentExpensesOutput
// @Failure     500
// @Router      /recurrent_expenses/all [get]
func (c *RecurrentExpensesController) getAll(ctx echo.Context) error {
	res, err := c.getAllRecurrentExpenses.GetAll(ctx.Request().Context())
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, res)
}

// @Summary     Create expenses from all recurrent expenses into db
// @Description Register expenses from all recurrent expenses for current month
// @Tags        recurrent expenses
// @Produce     json
// @Success     200 {object} recurrentexpenses.CreateMonthlyRecurrentExpensesOutput
// @Failure     500
// @Router      /recurrent_expenses/monthly_expenses [post]
func (c *RecurrentExpensesController) createMonthly(ctx echo.Context) error {
	got, err := c.createMonthlyRecurrentExpenses.CreateMonthlyRecurrentExpenses(ctx.Request().Context())
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, got)
}
