package expenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
)

type (
	CreateExpenseInput struct {
		Name        string
		Amount      float64
		Description string
	}
	CreateExpense interface {
		Create(context.Context, *CreateExpenseInput) (*entities.Expense, error)
	}
	CreateExpenseImpl struct {
		expensesRepo repos.ExpensesRepository
	}
)

func NewCreateExpensesImpl(repo repos.ExpensesRepository) *CreateExpenseImpl {
	return &CreateExpenseImpl{
		expensesRepo: repo,
	}
}

func (c *CreateExpenseImpl) Create(ctx context.Context, expense *CreateExpenseInput) (*entities.Expense, error) {
	newExpense := entities.Expense{
		Name:        expense.Name,
		Amount:      expense.Amount,
		Description: expense.Description,
	}
	if err := c.expensesRepo.Save(ctx, &newExpense); err != nil {
		return nil, err
	}
	return &newExpense, nil
}
