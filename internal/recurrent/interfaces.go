package recurrent

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	RecurrentExpenseCreator interface {
		Create(ctx context.Context, recurrentExpense *entities.RecurrentExpense) error
	}
	RecurrentExpensesCounter interface {
		// CountRecurrentExpensesByDateAndID count all recurrent expenses in existance from date starting and date ending
		CountRecurrentExpensesByDateAndID(date time.Time, recurrentExpenseID uuid.UUID) (int64, error)
	}
	RecurrentExpensesGetter interface {
		GetRecurrentExpensesByDate(date time.Time) ([]entities.RecurrentExpense, error)
	}
)
