package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/middlewares"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/json"
)

type RecurrentExpensesController struct {
	middlewares.Middlewares
	createRecurrentExpense         recurrentexpenses.RecurrentExpenseCreatable
	getAllRecurrentExpenses        recurrentexpenses.RecurrentExpensesAllGettable
	createMonthlyRecurrentExpenses recurrentexpenses.MonthlyRecurrentExpensesCreateable
	group                          *echo.Group
}

func NewRecurrentExpensesController(
	createRecurrentExpense recurrentexpenses.RecurrentExpenseCreatable,
	getAllRecurrentExpenses recurrentexpenses.RecurrentExpensesAllGettable,
	createMonthlyRecurrentExpenses recurrentexpenses.MonthlyRecurrentExpensesCreateable,
	middlewares middlewares.Middlewares,
	e *echo.Echo, //nolint:varnamelen
) *RecurrentExpensesController {
	controller := &RecurrentExpensesController{
		Middlewares:                    middlewares,
		createRecurrentExpense:         createRecurrentExpense,
		getAllRecurrentExpenses:        getAllRecurrentExpenses,
		createMonthlyRecurrentExpenses: createMonthlyRecurrentExpenses,
		group:                          e.Group("/recurrent_expenses"),
	}
	controller.register()
	return controller
}

func (c *RecurrentExpensesController) register() {
	c.group.POST("", c.Create, c.LoggedIn)
	c.group.GET("/all", c.GetAll, c.LoggedIn)
	c.group.POST("/monthly_expenses", c.CreateMonthly, c.LoggedIn)
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
// @Security    ApiKeyAuth
// @Router      /recurrent_expenses [post]
func (c *RecurrentExpensesController) Create(ctx echo.Context) error {
	var request recurrentexpenses.CreateRecurrentExpenseInput
	if err := ctx.Bind(&request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	request.UserID = ctx.Get("user_id").(string)
	if err := ctx.Validate(request); err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	log.Infoln(json.MustMarshall(request))
	res, err := c.createRecurrentExpense.CreateRecurrentExpense(ctx.Request().Context(), &request)
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
// @Security    ApiKeyAuth
// @Router      /recurrent_expenses/all [get]
func (c *RecurrentExpensesController) GetAll(ctx echo.Context) error {
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
// @Security    ApiKeyAuth
// @Router      /recurrent_expenses/monthly_expenses [post]
func (c *RecurrentExpensesController) CreateMonthly(ctx echo.Context) error {
	got, err := c.createMonthlyRecurrentExpenses.CreateMonthlyRecurrentExpenses(ctx.Request().Context(), uuid.MustParse(ctx.Get("user_id").(string)))
	if err != nil {
		return apperrors.CreateResponseFromError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, got)
}
