package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/pkg/errors"
)

type ExpensesController struct {
	createExpense   expenses.ExpenseCreatable
	setToPaid       expenses.ExpenseToPaidSetteable
	togglableIsPaid expenses.ExpenseToPaidTogglable
	updateExpense   expenses.ExpenseUpdateable
	group           *echo.Group
}

func NewExpensesController(
	createExpense expenses.ExpenseCreatable,
	setToPaid expenses.ExpenseToPaidSetteable,
	togglableIsPaid expenses.ExpenseToPaidTogglable,
	updateExpense expenses.ExpenseUpdateable,
	e *echo.Echo, //nolint:varnamelen
) *ExpensesController {
	return &ExpensesController{
		createExpense:   createExpense,
		setToPaid:       setToPaid,
		togglableIsPaid: togglableIsPaid,
		updateExpense:   updateExpense,
		group:           e.Group("/expenses"),
	}
}

func (c *ExpensesController) Register() {
	c.group.POST("", c.create)
	c.group.POST("/to_paid", c.toPaid)
	c.group.POST("/toggle_is_paid", c.toggleIsPaid)
	c.group.PUT("/update", c.update)
}

// @Summary     Create an expense
// @Description Register a expense into the database
// @Tags        expenses
// @Accept      json
// @Produce     json
// @Param       expense_to_create body     expenses.CreateExpenseInput true "Expense to be created"
// @Success     201               {object} entities.Expense            "Expense has been created"
// @Failure     400               {object} validator.ValidationError   "When a request does not fulfill need data"
// @Failure     500               "Something unidentified has occurred"
// @Router      /expenses [post]
func (c *ExpensesController) create(ctx echo.Context) error {
	var expenseRequest expenses.CreateExpenseInput
	if err := ctx.Bind(&expenseRequest); err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	newExpense, err := c.createExpense.CreateExpense(ctx.Request().Context(), &expenseRequest)
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
// @Failure     400 {object} validator.ValidationError "When a request does not fulfill need data"
// @Failure     500
// @Router      /expenses/to_paid [post]
func (c *ExpensesController) toPaid(ctx echo.Context) error {
	var request expenses.SetExpenseToPaidInput
	if err := ctx.Bind(&request); err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	if err := c.setToPaid.SetToPaid(ctx.Request().Context(), &request); err != nil { //nolint: staticcheck
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.NoContent(http.StatusOK)
}

// @Summary     Toggle is paid status from an expense
// @Description Toggle is paid status from an expense
// @Tags        expenses
// @Accept      json
// @Produce     json
// @Param       expense_id body     expenses.ToggleExpenseIsPaidInput true "ID to toggle is paid status"
// @Success     200        {object} expenses.ToggleExpenseIsPaidOutput
// @Failure     400        {object} validator.ValidationError "When a request does not fulfill need data"
// @Failure     500
// @Router      /expenses/toggle_is_paid [post]
func (c *ExpensesController) toggleIsPaid(ctx echo.Context) error {
	var request expenses.ToggleExpenseIsPaidInput
	if err := ctx.Bind(&request); err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	got, err := c.togglableIsPaid.ToggleIsPaid(ctx.Request().Context(), &request)
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, got)
}

// @Summary     Updates an expense
// @Description Updates an expense data depending if it is recurrent or not
// @Tags        expenses
// @Accept      json
// @Produce     json
// @Param       expense_update_data body expenses.UpdateExpenseInput true "ID to toggle is paid status"
// @Success     200                 "Expense was updated"
// @Failure     400                 {object} validator.ValidationError "When a request does not fulfill need data"
// @Failure     500
// @Router      /expenses/update [put]
func (c *ExpensesController) update(ctx echo.Context) error {
	var request expenses.UpdateExpenseInput
	if err := ctx.Bind(&request); err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	err := c.updateExpense.UpdateExpense(ctx.Request().Context(), &request)
	if err != nil {
		return errors.CreateResponseFromError(ctx, err)
	}
	return ctx.NoContent(http.StatusOK)
}
