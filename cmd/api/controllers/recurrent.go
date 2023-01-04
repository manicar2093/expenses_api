package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/middlewares"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
)

type RecurrentExpensesController struct {
	*middlewares.EchoMiddlewares
	createRecurrentExpense         recurrentexpenses.RecurrentExpenseCreatable
	getAllRecurrentExpenses        recurrentexpenses.RecurrentExpensesAllGettable
	createMonthlyRecurrentExpenses recurrentexpenses.MonthlyRecurrentExpensesCreateable
	group                          *echo.Group
}

func NewRecurrentExpensesController(
	createRecurrentExpense recurrentexpenses.RecurrentExpenseCreatable,
	getAllRecurrentExpenses recurrentexpenses.RecurrentExpensesAllGettable,
	createMonthlyRecurrentExpenses recurrentexpenses.MonthlyRecurrentExpensesCreateable,
	middlewares *middlewares.EchoMiddlewares,
	e *echo.Echo, //nolint:varnamelen
) *RecurrentExpensesController {
	return &RecurrentExpensesController{
		EchoMiddlewares:                middlewares,
		createRecurrentExpense:         createRecurrentExpense,
		getAllRecurrentExpenses:        getAllRecurrentExpenses,
		createMonthlyRecurrentExpenses: createMonthlyRecurrentExpenses,
		group:                          e.Group("/recurrent_expenses"),
	}
}

func (c *RecurrentExpensesController) Register() {
	c.group.POST("", c.create, c.LoggedIn)

	c.group.GET("/all", c.getAll, c.LoggedIn)

	c.group.POST("/monthly_expenses", c.createMonthly, c.LoggedIn)
}

// @Summary     Create a recurrent expense
// @Description Register into db a recurrent expense
// @Tags        recurrent expenses
// @Accept      json
// @Produce     json
// @Param       recurrent_expense body recurrentexpenses.CreateRecurrentExpenseInput true "Recurrent Expense"
// @Success     201
// @Failure     400 {object} validator.ValidationError "When a request does not fulfill need data"
// @Failure     500
// @Router      /recurrent_expenses [post]
func (c *RecurrentExpensesController) create(ctx echo.Context) error {
	var recurrentExpenseReq recurrentexpenses.CreateRecurrentExpenseInput
	if err := ctx.Bind(&recurrentExpenseReq); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	recurrentExpenseReq.UserID = ctx.Get("user_id").(string)
	res, err := c.createRecurrentExpense.CreateRecurrentExpense(ctx.Request().Context(), &recurrentExpenseReq)
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, res)
}

// @Summary     Get all recurrent expense
// @Description Finds all recurrent expenses registered in db
// @Tags        recurrent expenses
// @Produce     json
// @Success     200
// @Failure     500
// @Router      /recurrent_expenses/all [get]
func (c *RecurrentExpensesController) getAll(ctx echo.Context) error {
	res, err := c.getAllRecurrentExpenses.GetAll(ctx.Request().Context(), uuid.MustParse(ctx.Get("user_id").(string)))
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, res)
}

// @Summary     Create expenses from all recurrent expenses into db
// @Description Register expenses from all recurrent expenses for current month
// @Tags        recurrent expenses
// @Produce     json
// @Success     200
// @Failure     500
// @Router      /recurrent_expenses/monthly_expenses [post]
func (c *RecurrentExpensesController) createMonthly(ctx echo.Context) error {
	got, err := c.createMonthlyRecurrentExpenses.CreateMonthlyRecurrentExpenses(ctx.Request().Context(), uuid.MustParse(ctx.Get("user_id").(string)))
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, got)
}
