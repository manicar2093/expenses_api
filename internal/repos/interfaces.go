package repos

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	ExpensesRepository interface {
		Save(ctx context.Context, expense *entities.Expense) error
		UpdateIsPaidByExpenseID(ctx context.Context, expenseID uuid.UUID, status bool) error
		GetExpenseStatusByID(ctx context.Context, expenseID uuid.UUID) (*entities.ExpenseIDWithIsPaidStatus, error)
		Update(ctx context.Context, expenseUpdateInput *UpdateExpenseInput) error
		FindByID(ctx context.Context, expenseID uuid.UUID) (*entities.Expense, error)
	}

	IncomesRepository interface {
		Save(context.Context, *entities.Income) error
	}

	RecurrentExpenseRepo interface {
		Save(ctx context.Context, recExpense *entities.RecurrentExpense) error
		FindByName(ctx context.Context, name string, userID uuid.UUID) (*entities.RecurrentExpense, error)
		FindAll(ctx context.Context, userID uuid.UUID) ([]*entities.RecurrentExpense, error)
	}

	UserRepo interface {
		Save(ctx context.Context, user *entities.User) error
		FindUserByEmail(ctx context.Context, email string) (*entities.User, error)
	}
)
