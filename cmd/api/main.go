package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manicar2093/expenses_api/internal/connections"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/internal/incomes"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/errors"
)

var (
	mongoConn             = connections.GetMongoConn()
	expensesRepo          = repos.NewExpensesRepositoryImpl(mongoConn)
	recurrentExpensesRepo = repos.NewRecurrentExpenseRepoImpl(mongoConn)
	timeGetter            = &dates.TimeGetter{}
	e                     = echo.New() //nolint:varnamelen
)

func main() {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	expensesRoutes()
	incomesRouter()
	reportsRoutes()
	recurrentExpensesRoutes()
	e.Logger.Fatal(e.Start(":8000"))
}

func incomesRouter() {
	var (
		incomesRepo  = repos.NewIncomesRepositoryImpl(mongoConn)
		createIncome = incomes.NewCreateIncomeImpl(incomesRepo)
		incomesGroup = e.Group("/incomes")
	)
	incomesGroup.POST("", func(ctx echo.Context) error {
		var incomeRequest incomes.CreateIncomeInput
		if err := ctx.Bind(&incomeRequest); err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		newIncome, err := createIncome.Create(ctx.Request().Context(), &incomeRequest)
		if err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		return ctx.JSON(http.StatusCreated, newIncome)
	})
}

func expensesRoutes() {
	var (
		createExpense    = expenses.NewCreateExpensesImpl(expensesRepo, timeGetter)
		setExpenseToPaid = expenses.NewSetExpenseToPaidImpl(expensesRepo)
		expensesGroup    = e.Group("/expenses")
	)
	expensesGroup.POST("", func(ctx echo.Context) error {
		var expenseRequest expenses.CreateExpenseInput
		if err := ctx.Bind(&expenseRequest); err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		newExpense, err := createExpense.Create(ctx.Request().Context(), &expenseRequest)
		if err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		return ctx.JSON(http.StatusCreated, newExpense)
	})

	expensesGroup.POST("/to_paid", func(ctx echo.Context) error {
		var request expenses.SetExpenseToPaidInput
		if err := ctx.Bind(&request); err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		if err := setExpenseToPaid.SetToPaid(ctx.Request().Context(), &request); err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		return ctx.NoContent(http.StatusOK)
	})
}

func reportsRoutes() {
	var (
		getCurrentMonth = reports.NewCurrentMonthDetailsImpl(expensesRepo, timeGetter)
		reportsGroup    = e.Group("/reports")
	)

	reportsGroup.GET("/current_month", func(ctx echo.Context) error {
		currentMonthDetails, err := getCurrentMonth.GetExpenses(ctx.Request().Context())
		if err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		return ctx.JSON(http.StatusOK, currentMonthDetails)
	})
}

func recurrentExpensesRoutes() {
	var (
		createRecurrentExpense = recurrentexpenses.NewCreateRecurrentExpenseImpl(
			recurrentExpensesRepo,
			expensesRepo,
			timeGetter,
		)
		getAllRecurrentExpenses = recurrentexpenses.NewGetAllRecurrentExpensesImpl(recurrentExpensesRepo)
		recurrentExpenseGroup   = e.Group("/recurrent_expenses")
	)

	recurrentExpenseGroup.POST("", func(ctx echo.Context) error {
		var recurrentExpenseReq recurrentexpenses.CreateRecurrentExpenseInput
		if err := ctx.Bind(&recurrentExpenseReq); err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		res, err := createRecurrentExpense.Create(ctx.Request().Context(), &recurrentExpenseReq)
		if err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		return ctx.JSON(http.StatusCreated, res)
	})

	recurrentExpenseGroup.GET("/all", func(ctx echo.Context) error {
		res, err := getAllRecurrentExpenses.GetAll(ctx.Request().Context())
		if err != nil {
			return errors.CreateResponseFromError(ctx, err)
		}
		return ctx.JSON(http.StatusOK, res)
	})
}
