package repos

import (
	"context"
	"time"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/schemas"
)

type (
	ExpensesRepository interface {
		Save(ctx context.Context, expense *mongoentities.Expense) error
		GetExpensesByMonth(ctx context.Context, month time.Month) ([]*mongoentities.Expense, error)
		UpdateIsPaidByExpenseID(ctx context.Context, expenseID interface{}, status bool) error
		FindByNameAndMonthAndIsRecurrent(ctx context.Context, month uint, expenseName string) (*mongoentities.Expense, error)
		GetExpenseStatusByID(ctx context.Context, expenseID interface{}) (*schemas.ExpenseIDWithIsPaidStatus, error)
	}

	IncomesRepository interface {
		Save(context.Context, *mongoentities.Income) error
	}

	RecurrentExpenseRepo interface {
		Save(ctx context.Context, recExpense *mongoentities.RecurrentExpense) error
		FindByName(ctx context.Context, name string) (*mongoentities.RecurrentExpense, error)
		FindAll(ctx context.Context) (*[]mongoentities.RecurrentExpense, error)
	}
)
