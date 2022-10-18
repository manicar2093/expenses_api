package expenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	ExpenseToPaidSetteable interface {
		// Deprecated: Use ToggleIsPaid instead
		SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error
	}
	ExpenseCreatable interface {
		Create(context.Context, *CreateExpenseInput) (*entities.Expense, error)
	}
	ExpenseToPaidTogglable interface {
		ToggleIsPaid(ctx context.Context, input *ToggleExpenseIsPaidInput) (*ToggleExpenseIsPaidOutput, error)
	}
	ExpenseSevice interface {
		ExpenseToPaidSetteable
		ExpenseCreatable
		ExpenseToPaidTogglable
	}
)
