package incomes

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/json"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

type (
	CreateIncome interface {
		Create(ctx context.Context, incomeInput *CreateIncomeInput) (*entities.Income, error)
	}

	IncomeServiceImpl struct {
		incomesRepo repos.IncomesRepository
		validator   validator.StructValidable
	}
)

func NewIncomeServiceImpl(repo repos.IncomesRepository, validator validator.StructValidable) *IncomeServiceImpl {
	return &IncomeServiceImpl{
		incomesRepo: repo,
		validator:   validator,
	}
}

func (c *IncomeServiceImpl) Create(ctx context.Context, incomeInput *CreateIncomeInput) (*entities.Income, error) {
	log.Infoln("Request: ", json.MustMarshall(incomeInput))
	if err := c.validator.ValidateStruct(incomeInput); err != nil {
		return nil, err
	}

	if err := c.incomesRepo.Save(ctx, &incomeInput.Income); err != nil {
		return nil, err
	}
	return &incomeInput.Income, nil
}
