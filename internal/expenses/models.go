package expenses

import (
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

type (
	ExpenseServiceImpl struct {
		expensesRepo repos.ExpensesRepository
		timeGetter   dates.TimeGetable
		validator    validator.StructValidable
	}
)

func NewExpenseServiceImpl(repo repos.ExpensesRepository, timeGetter dates.TimeGetable, validator validator.StructValidable) *ExpenseServiceImpl {
	return &ExpenseServiceImpl{
		expensesRepo: repo,
		timeGetter:   timeGetter,
		validator:    validator,
	}
}
