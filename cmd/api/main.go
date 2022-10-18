package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manicar2093/expenses_api/cmd/api/controllers"
	_ "github.com/manicar2093/expenses_api/cmd/api/docs"
	"github.com/manicar2093/expenses_api/internal/connections"
	_ "github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	mongoConn    = connections.GetMongoConn()
	expensesRepo = repos.NewExpensesRepositoryImpl(
		mongoConn,
	)
	recurrentExpensesRepo = repos.NewRecurrentExpenseRepoImpl(
		mongoConn,
	)
	timeGetter     = &dates.TimeGetter{}
	expenseService = expenses.NewExpenseServiceImpl(
		expensesRepo,
		timeGetter,
	)
	getCurrentMonth = reports.NewCurrentMonthDetailsImpl(
		expensesRepo,
		timeGetter,
	)
	createRecurrentExpense = recurrentexpenses.NewCreateRecurrentExpenseImpl(
		recurrentExpensesRepo,
		expensesRepo,
		timeGetter,
	)
	getAllRecurrentExpenses = recurrentexpenses.NewGetAllRecurrentExpensesImpl(
		recurrentExpensesRepo,
	)
	createMonthlyRecurrentExpenses = recurrentexpenses.NewCreateMonthlyRecurrentExpensesImpl(
		recurrentExpensesRepo,
		expensesRepo,
		timeGetter,
	)
	e = echo.New() //nolint:varnamelen
)

// @title   Expenses API
// @version 1.0
func main() {
	configEcho()
	registerControllers()
	e.Logger.Fatal(e.Start(":8000"))
}

func configEcho() {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
}

func registerControllers() {
	controllers.NewExpensesController(
		expenseService,
		expenseService,
		expenseService,
		e,
	).Register()
	controllers.NewRecurrentExpensesController(
		createRecurrentExpense,
		getAllRecurrentExpenses,
		createMonthlyRecurrentExpenses,
		e,
	).Register()
	controllers.NewReportsController(
		getCurrentMonth,
		e,
	).Register()
	controllers.NewHealthCheckController(
		mongoConn.Client(),
		e,
	).Register()
}
