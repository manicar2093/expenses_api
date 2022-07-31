package recurrentexpenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
)

type (
	GetAllRecurrentExpenses interface {
		GetAll(ctx context.Context) (*[]entities.RecurrentExpense, error)
	}
	GetAllRecurrentExpensesImpl struct {
		recurrentExpensesRepo repos.RecurrentExpenseRepo
	}
)

func NewGetAllRecurrentExpensesImpl(recurrentExpensesRepo repos.RecurrentExpenseRepo) *GetAllRecurrentExpensesImpl {
	return &GetAllRecurrentExpensesImpl{recurrentExpensesRepo: recurrentExpensesRepo}
}

func (c *GetAllRecurrentExpensesImpl) GetAll(ctx context.Context) (*[]entities.RecurrentExpense, error) {
	return c.recurrentExpensesRepo.FindAll(ctx)
}
