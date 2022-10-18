package expenses

import (
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	SetExpenseToPaidInput struct {
		ID string `json:"id,omitempty"`
	}
	CreateExpenseInput struct {
		Name         string  `json:"name,omitempty"`
		Amount       float64 `json:"amount,omitempty"`
		Description  string  `json:"description,omitempty"`
		ForNextMonth bool    `json:"for_next_month,omitempty"`
	}
	ExpenseServiceImpl struct {
		expensesRepo repos.ExpensesRepository
		timeGetter   dates.TimeGetable
	}
)

func NewExpenseServiceImpl(repo repos.ExpensesRepository, timeGetter dates.TimeGetable) *ExpenseServiceImpl {
	return &ExpenseServiceImpl{
		expensesRepo: repo,
		timeGetter:   timeGetter,
	}
}
