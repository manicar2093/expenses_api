package expenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/json"
)

type (
	CreateExpenseInput struct {
		Name        string  `json:"name,omitempty"`
		Amount      float64 `json:"amount,omitempty"`
		Description string  `json:"description,omitempty"`
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
	log.Println("Request: ", json.MustMarshall(expense))
	newExpense := entities.Expense{
		Name:        expense.Name,
		Amount:      expense.Amount,
		Description: expense.Description,
		IsPaid:      true,
	}
	if err := c.expensesRepo.Save(ctx, &newExpense); err != nil {
		return nil, err
	}
	return &newExpense, nil
}
