package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/pkg/errors"
)

type ExpensesController struct {
	createExpense expenses.ExpenseCreatable
	setToPaid     expenses.ExpenseToPaidSetteable
	group         *echo.Group
}

func NewExpensesController(
	createExpense expenses.ExpenseCreatable,
	setToPaid expenses.ExpenseToPaidSetteable,
	e *echo.Echo,
) *ExpensesController {
	return &ExpensesController{
		createExpense: createExpense,
		setToPaid:     setToPaid,
		group:         e.Group("/expenses"),
	}
}

func (c *ExpensesController) Register() {
	c.group.POST("", c.create)
	c.group.POST("/to_paid", c.toPaid)
}

// @Summary     Create an expense
// @Description Register a expense into the database
// @Tags        expenses
// @Accept      json
// @Produce     json
// @Param       create_expense body     expenses.CreateExpenseInput true "Expense to be created"
// @Success     201            {object} entities.Expense
// @Failure     500
// @Router      /expenses [post]
func (c *ExpensesController) create(ctx echo.Context) error {
	var expenseRequest expenses.CreateExpenseInput
	if err := ctx.Bind(&expenseRequest); err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	newExpense, err := c.createExpense.Create(ctx.Request().Context(), &expenseRequest)
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, newExpense)
}

// @Summary     Set expense to paid
// @Description Change a expense is_paid status
// @Tags        expenses
// @Accept      json
// @Produce     json
// @Param       expense_id body expenses.SetExpenseToPaidInput true "ID to change to is paid"
// @Success     200
// @Failure     500
// @Router      /expenses/to_paid [post]
func (c *ExpensesController) toPaid(ctx echo.Context) error {
	var request expenses.SetExpenseToPaidInput
	if err := ctx.Bind(&request); err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	if err := c.setToPaid.SetToPaid(ctx.Request().Context(), &request); err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.NoContent(http.StatusOK)
}
