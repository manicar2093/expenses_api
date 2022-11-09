package periodizer

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	ExpensesCountByRecurrentExpensePeriodicityGenerable interface {
		GenerateExpensesCountByRecurrentExpensePeriodicity(
			ctx context.Context,
			recurrentExpense *entities.RecurrentExpense,
		) ([]*entities.Expense, error)
	}
)
