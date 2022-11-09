package periodicity

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ExpensePeriodicityAddable interface {
		AddRecurrentExpense(
			ctx context.Context,
			month, year uint,
			recurrentExpense *entities.RecurrentExpense,
		) (
			*entities.RecurrentExpensesMonthlyCreated,
			error,
		)
	}

	ExpensePeriodicityIsPaidToggable interface {
		ToggleExpenseIsPaidStatus(
			ctx context.Context,
			recurrentExpenseID, expenseID primitive.ObjectID,
			newIsPaidStatus bool,
		) (
			*entities.RecurrentExpensesMonthlyCreated,
			error,
		)
	}

	ExpensePeriodicityCreateable interface {
		GenerateRecurrentExpensesByYearAndMonth(
			ctx context.Context,
			month, year uint,
		) (
			*entities.RecurrentExpensesMonthlyCreated,
			error,
		)
	}

	ExpensesPeriodicityService interface {
		ExpensePeriodicityCreateable
	}
)
