package periodicity

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
)

type ExpensePeriodicityCreateable interface {
	GenerateRecurrentExpensesByYearAndMonth(
		ctx context.Context,
		month, year uint,
	) (
		*entities.RecurrentExpensesMonthlyCreated,
		error,
	)
}

