package recurrentexpenses

import (
	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	GetAllRecurrentExpensesOutput struct {
		RecurrentExpenses       []*entities.RecurrentExpense `json:"recurrent_expenses,omitempty"`
		RecurrenteExpensesCount uint                         `json:"recurrente_expenses_count,omitempty"`
	}
)
