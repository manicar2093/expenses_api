package expenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	CurrentMonthDetailsOutput struct {
		TotalAmount   float64
		TotalExpenses uint
		Expenses      *[]entities.Expense
	}
	CurrentMonthDetailsGettable interface {
		GetExpenses(ctx context.Context) (*CurrentMonthDetailsOutput, error)
	}
	CurrentMonthDetails struct {
		repo       repos.ExpensesRepository
		timeGetter dates.TimeGetable
	}
)

func NewCurrentMonthDetailsImpl(repo repos.ExpensesRepository, timeGetter dates.TimeGetable) *CurrentMonthDetails {
	return &CurrentMonthDetails{repo: repo, timeGetter: timeGetter}
}

func (c *CurrentMonthDetails) GetExpenses(ctx context.Context) (*CurrentMonthDetailsOutput, error) {
	currentExpenses, err := c.repo.GetExpensesByMonth(ctx, c.timeGetter.GetCurrentTime().Month())
	if err != nil {
		return nil, err
	}

	return &CurrentMonthDetailsOutput{
		TotalAmount:   c.calculateTotalAmount(currentExpenses),
		TotalExpenses: uint(len(*currentExpenses)),
		Expenses:      currentExpenses,
	}, nil

}

func (c *CurrentMonthDetails) calculateTotalAmount(expenses *[]entities.Expense) float64 {
	var total float64
	for _, v := range *expenses {
		total += v.Amount
	}
	return total
}
