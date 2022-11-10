package incomes

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/json"
)

type (
	CreateIncome interface {
		Create(ctx context.Context, incomeInput *CreateIncomeInput) (*mongoentities.Income, error)
	}
	CreateIncomeInput struct {
		Name        string  `json:"name,omitempty"`
		Amount      float64 `json:"amount,omitempty"`
		Description string  `json:"description,omitempty"`
	}
	CreateIncomeImpl struct {
		incomesRepo repos.IncomesRepository
	}
)

func NewCreateIncomeImpl(repo repos.IncomesRepository) *CreateIncomeImpl {
	return &CreateIncomeImpl{
		incomesRepo: repo,
	}
}

func (c *CreateIncomeImpl) Create(ctx context.Context, incomeInput *CreateIncomeInput) (*mongoentities.Income, error) {
	log.Println("Request: ", json.MustMarshall(incomeInput))
	newIncome := mongoentities.Income{
		Name:        incomeInput.Name,
		Amount:      incomeInput.Amount,
		Description: incomeInput.Description,
	}
	if err := c.incomesRepo.Save(ctx, &newIncome); err != nil {
		return nil, err
	}
	return &newIncome, nil
}
