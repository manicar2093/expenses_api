package recurrentexpenses

import "github.com/manicar2093/expenses_api/internal/entities"

type (
	CreateRecurrentExpenseInput struct {
		Name        string  `json:"name,omitempty" validate:"required"`
		Amount      float64 `json:"amount,omitempty" validate:"required"`
		Description string  `json:"description,omitempty" validate:"-"`
		UserID      string  `json:"-"`
	}
	CreateRecurrentExpenseOutput struct {
		RecurrentExpense *entities.RecurrentExpense `json:"recurrent_expense,omitempty"`
		NextMonthExpense *entities.Expense          `json:"next_month_expense,omitempty"`
	}
)
