package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	ExpensesRepository interface {
		Save(ctx context.Context, expense *entities.Expense) error
	}
)
