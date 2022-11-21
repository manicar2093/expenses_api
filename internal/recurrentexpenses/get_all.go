package recurrentexpenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
)

type (
	GetAllRecurrentExpenses interface {
		GetAll(ctx context.Context) (*GetAllRecurrentExpensesOutput, error)
	}

	GetAllRecurrentExpensesOutput struct {
		RecurrentExpenses       []*entities.RecurrentExpense `json:"recurrent_expenses,omitempty"`
		RecurrenteExpensesCount uint                         `json:"recurrente_expenses_count,omitempty"`
	}
	GetAllRecurrentExpensesImpl struct {
		recurrentExpensesRepo repos.RecurrentExpenseRepo
	}
)

func NewGetAllRecurrentExpensesImpl(recurrentExpensesRepo repos.RecurrentExpenseRepo) *GetAllRecurrentExpensesImpl {
	return &GetAllRecurrentExpensesImpl{recurrentExpensesRepo: recurrentExpensesRepo}
}

func (c *GetAllRecurrentExpensesImpl) GetAll(ctx context.Context) (*GetAllRecurrentExpensesOutput, error) {
	recurrentExpenses, err := c.recurrentExpensesRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return &GetAllRecurrentExpensesOutput{
		RecurrentExpenses:       recurrentExpenses,
		RecurrenteExpensesCount: uint(len(recurrentExpenses)),
	}, nil
}
