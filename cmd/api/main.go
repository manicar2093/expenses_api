package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/manicar2093/expenses_api/internal/connections"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/internal/repos"
)

var (
	conn          = connections.GetRelConnection()
	expensesRepo  = repos.NewExpensesRepositoryImpl(conn)
	createExpense = expenses.NewCreateExpensesImpl(expensesRepo)
	e             = echo.New()
)

func main() {
	expensesRoutes()
	e.Logger.Fatal(e.Start(":8000"))
}

func expensesRoutes() {
	expensesGroup := e.Group("/expenses")
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
