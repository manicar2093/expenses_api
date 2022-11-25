package repos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	ExpensesRepository interface {
		Save(ctx context.Context, expense *entities.Expense) error
		GetExpensesByMonth(ctx context.Context, month time.Month) ([]*entities.Expense, error)
		UpdateIsPaidByExpenseID(ctx context.Context, expenseID uuid.UUID, status bool) error
		FindByNameAndMonthAndIsRecurrent(ctx context.Context, month uint, expenseName string) (*entities.Expense, error)
		GetExpenseStatusByID(ctx context.Context, expenseID uuid.UUID) (*entities.ExpenseIDWithIsPaidStatus, error)
		Update(ctx context.Context, expenseUpdateInput *UpdateExpenseInput) error
	}

	IncomesRepository interface {
		Save(context.Context, *entities.Income) error
	}

	RecurrentExpenseRepo interface {
		Save(ctx context.Context, recExpense *entities.RecurrentExpense) error
		FindByName(ctx context.Context, name string) (*entities.RecurrentExpense, error)
		FindAll(ctx context.Context) ([]*entities.RecurrentExpense, error)
	}
)
