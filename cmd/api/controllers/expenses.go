package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/middlewares"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/json"
)

type ExpensesController struct {
	middlewares.Middlewares
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
	middlewares middlewares.Middlewares,
	e *echo.Echo, //nolint:varnamelen
) *ExpensesController {
	controller := &ExpensesController{
		Middlewares:     middlewares,
		createExpense:   createExpense,
		setToPaid:       setToPaid,
		togglableIsPaid: togglableIsPaid,
		updateExpense:   updateExpense,
		group:           e.Group("/expenses"),
	}
	controller.register()
	return controller
}

func (c *ExpensesController) register() {
	c.group.POST("", c.Create, c.LoggedIn)
	c.group.PUT("/to_paid", c.ToPaid, c.LoggedIn)
	c.group.PUT("/toggle_is_paid", c.ToggleIsPaid, c.LoggedIn)
	c.group.PUT("/update", c.Update, c.LoggedIn)
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
// @Security    ApiKeyAuth
// @Router      /expenses [post]
func (c *ExpensesController) Create(ctx echo.Context) error {
	var expenseRequest expenses.CreateExpenseInput
	if err := ctx.Bind(&expenseRequest); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	expenseRequest.UserID = ctx.Get("user_id").(string)
	if err := ctx.Validate(expenseRequest); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	newExpense, err := c.createExpense.CreateExpense(ctx.Request().Context(), &expenseRequest)
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
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
// @Security    ApiKeyAuth
// @Router      /expenses/to_paid [put]
func (c *ExpensesController) ToPaid(ctx echo.Context) error {
	log.Warnln("DEPRECATED!. Use ExpenseToPaidTogglable instead")
	var request expenses.SetExpenseToPaidInput
	if err := ctx.Bind(&request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	if err := ctx.Validate(request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	log.Infoln(json.MustMarshall(request))
	if err := c.setToPaid.SetToPaid(ctx.Request().Context(), &request); err != nil { //nolint: staticcheck
		return apperrors.CreateResponseFromError(ctx, err)
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
// @Security    ApiKeyAuth
// @Router      /expenses/toggle_is_paid [put]
func (c *ExpensesController) ToggleIsPaid(ctx echo.Context) error {
	var request expenses.ToggleExpenseIsPaidInput
	if err := ctx.Bind(&request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	if err := ctx.Validate(request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	log.Infoln(json.MustMarshall(request))
	got, err := c.togglableIsPaid.ToggleIsPaid(ctx.Request().Context(), &request)
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
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
// @Security    ApiKeyAuth
// @Router      /expenses/update [put]
func (c *ExpensesController) Update(ctx echo.Context) error {
	var request expenses.UpdateExpenseInput
	if err := ctx.Bind(&request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	if err := ctx.Validate(request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	log.Infoln(json.MustMarshall(request))
	err := c.updateExpense.UpdateExpense(ctx.Request().Context(), &request)
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	return ctx.NoContent(http.StatusOK)
}
