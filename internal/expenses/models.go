package expenses

import (
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

type (
	ToggleExpenseIsPaidInput struct {
		ID uuid.UUID `json:"id,omitempty"`
	}
	ToggleExpenseIsPaidOutput struct {
		ID                  uuid.UUID `json:"id"`
		CurrentIsPaidStatus bool      `json:"current_is_paid_status"`
	}
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
