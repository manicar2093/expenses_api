package expenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	ExpenseToPaidSetteable interface {
		SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error
	}
	ExpenseCreatable interface {
		Create(context.Context, *CreateExpenseInput) (*entities.Expense, error)
	}
	ExpenseSevice interface {
		ExpenseToPaidSetteable
		ExpenseCreatable
	}
)
