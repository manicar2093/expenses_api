package expenses

import (
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	ExpenseServiceImpl struct {
		expensesRepo repos.ExpensesRepository
		timaGettable dates.TimeGetable
	}
)

func NewExpenseServiceImpl(repo repos.ExpensesRepository, timeGetter dates.TimeGetable) *ExpenseServiceImpl {
	return &ExpenseServiceImpl{
		expensesRepo: repo,
		timaGettable: timeGetter,
	}
}
