package recurrentexpenses

import "context"

type (
	RecurrentExpenseCreatable interface {
		CreateRecurrentExpense(ctx context.Context, input *CreateRecurrentExpenseInput) (*CreateRecurrentExpenseOutput, error)
	}

	RecurrentExpensesAllGettable interface {
		GetAll(ctx context.Context) (*GetAllRecurrentExpensesOutput, error)
	}

	MonthlyRecurrentExpensesCreateable interface {
		CreateMonthlyRecurrentExpenses(ctx context.Context) (*CreateMonthlyRecurrentExpensesOutput, error)
	}
)
