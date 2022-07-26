package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/manicar2093/expenses_api/internal/connections"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/internal/incomes"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

var (
	mongoConn    = connections.GetMongoConn()
	expensesRepo = repos.NewExpensesRepositoryImpl(mongoConn)
	timeGetter   = dates.TimeGetter{}
	e            = echo.New() //nolint:varnamelen
)

func main() {
	expensesRoutes()
	incomesRouter()
	reportsRoutes()
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
			return err
		}
		newIncome, err := createIncome.Create(ctx.Request().Context(), &incomeRequest)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusCreated, newIncome)
	})
}

func expensesRoutes() {
	var (
		createExpense = expenses.NewCreateExpensesImpl(expensesRepo)
		expensesGroup = e.Group("/expenses")
	)
	expensesGroup.POST("", func(ctx echo.Context) error {
		var expenseRequest expenses.CreateExpenseInput
		if err := ctx.Bind(&expenseRequest); err != nil {
			return err
		}
		newExpense, err := createExpense.Create(ctx.Request().Context(), &expenseRequest)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusCreated, newExpense)
	})
}

func reportsRoutes() {
	var (
		getCurrentMonth = reports.NewCurrentMonthDetailsImpl(expensesRepo, &timeGetter)
		reportsGroup    = e.Group("/reports")
	)

	reportsGroup.GET("/current_month", func(ctx echo.Context) error {
		currentMonthDetails, err := getCurrentMonth.GetExpenses(ctx.Request().Context())
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, currentMonthDetails)
	})
}
