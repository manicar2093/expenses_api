package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/manicar2093/expenses_api/cmd/api/controllers"
	_ "github.com/manicar2093/expenses_api/cmd/api/docs"
	"github.com/manicar2093/expenses_api/cmd/api/middlewares"
	"github.com/manicar2093/expenses_api/internal/auth"
	"github.com/manicar2093/expenses_api/internal/config"
	"github.com/manicar2093/expenses_api/internal/connections"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/internal/tokens"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/validator"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	conn                  = connections.GetGormConnection()
	expensesRepo          = repos.NewExpensesGormRepo(conn)
	recurrentExpensesRepo = repos.NewRecurrentExpenseGormRepo(conn)
	timeGetter            = &dates.TimeGetter{}
	structValidator       = validator.NewGooKitValidator()
	customMiddlewares     = middlewares.NewEchoMiddlewares(tokens.NewPaseto(config.Instance.TokenSymmetricKey))
	expenseService        = expenses.NewExpenseServiceImpl(
		expensesRepo,
		timeGetter,
		structValidator,
	)
	getCurrentMonth = reports.NewCurrentMonthDetailsImpl(
		expensesRepo,
		timeGetter,
	)
	createRecurrentExpense = recurrentexpenses.NewCreateRecurrentExpense(
		recurrentExpensesRepo,
		expensesRepo,
		timeGetter,
		structValidator,
	)
	getAllRecurrentExpenses = recurrentexpenses.NewGetAllRecurrentExpenseServiceImpl(
		recurrentExpensesRepo,
	)
	createMonthlyRecurrentExpenses = recurrentexpenses.NewCreateMonthlyRecurrentExpensesImpl(
		recurrentExpensesRepo,
		expensesRepo,
		timeGetter,
	)
	e = echo.New() //nolint:varnamelen
)

// @title                      Expenses API
// @version                    1.0
// @securityDefinitions.apikey ApiKeyAuth
// @name                       Authorization
// @in                         header
// @authorizationurl           /auth/login/google
// @description                Type "Bearer" and then your API Token
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
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))
	e.Validator = structValidator
}

func registerControllers() {
	controllers.NewExpensesController(
		expenseService,
		expenseService,
		expenseService,
		expenseService,
		customMiddlewares,
		e,
	)
	controllers.NewRecurrentExpensesController(
		createRecurrentExpense,
		getAllRecurrentExpenses,
		createMonthlyRecurrentExpenses,
		customMiddlewares,
		e,
	)
	controllers.NewReportsController(
		getCurrentMonth,
		customMiddlewares,
		e,
	)
	controllers.NewHealthCheckController(
		conn,
		e,
	)
	controllers.NewLoginController(
		auth.NewGoogleTokenAuth(
			repos.NewUserGormRepo(conn),
			tokens.NewPaseto(config.Instance.TokenSymmetricKey),
			validator.NewGoogleTokenValidator(),
			config.Instance.AccessTokenDuration,
		),
		e,
	)
}
