package incomes

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/json"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

type (
	CreateIncome interface {
		Create(ctx context.Context, incomeInput *CreateIncomeInput) (*entities.Income, error)
	}
	CreateIncomeInput struct {
		Name        string    `json:"name,omitempty" validate:"required"`
		Amount      float64   `json:"amount,omitempty" validate:"required"`
		Description string    `json:"description,omitempty" validate:"-"`
		UserID      uuid.UUID `json:"user_uuid" validate:"required"`
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
	log.Println("Request: ", json.MustMarshall(incomeInput))
	if err := c.validator.ValidateStruct(incomeInput); err != nil {
		return nil, err
	}
	newIncome := entities.Income{
		Name:        incomeInput.Name,
		Amount:      incomeInput.Amount,
		Description: incomeInput.Description,
		UserID:      incomeInput.UserID,
	}
	if err := c.incomesRepo.Save(ctx, &newIncome); err != nil {
		return nil, err
	}
	return &newIncome, nil
}
