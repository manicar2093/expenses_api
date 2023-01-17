package incomes

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
)

type (
	CreateIncome interface {
		Create(ctx context.Context, incomeInput *CreateIncomeInput) (*entities.Income, error)
	}

	IncomeServiceImpl struct {
		incomesRepo repos.IncomesRepository
	}
)

func NewIncomeServiceImpl(repo repos.IncomesRepository) *IncomeServiceImpl {
	return &IncomeServiceImpl{
		incomesRepo: repo,
	}
}

func (c *IncomeServiceImpl) Create(ctx context.Context, incomeInput *CreateIncomeInput) (*entities.Income, error) {
	if err := c.incomesRepo.Save(ctx, &incomeInput.Income); err != nil {
		return nil, err
	}
	return &incomeInput.Income, nil
}
