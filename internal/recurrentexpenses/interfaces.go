package recurrentexpenses

import (
	"context"

	"github.com/google/uuid"
)

type (
	RecurrentExpenseCreatable interface {
		CreateRecurrentExpense(ctx context.Context, input *CreateRecurrentExpenseInput) (*CreateRecurrentExpenseOutput, error)
	}

	RecurrentExpensesAllGettable interface {
		GetAll(ctx context.Context, userID uuid.UUID) (*GetAllRecurrentExpensesOutput, error)
	}

	MonthlyRecurrentExpensesCreateable interface {
		CreateMonthlyRecurrentExpenses(ctx context.Context, userID uuid.UUID) (*CreateMonthlyRecurrentExpensesOutput, error)
	}
)
