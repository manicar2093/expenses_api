package recurrentexpenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/json"
)

type (
	CreateRecurrentExpense interface {
		Create(ctx context.Context, input *CreateRecurrentExpenseInput) (*CreateRecurrentExpenseOutput, error)
	}
	CreateRecurrentExpenseInput struct {
		Name        string  `json:"name,omitempty"`
		Amount      float64 `json:"amount,omitempty"`
		Description string  `json:"description,omitempty"`
	}
	CreateRecurrentExpenseOutput struct {
		RecurrentExpense *entities.RecurrentExpense `json:"recurrent_expense,omitempty"`
	}
	CreateRecurrentExpenseImpl struct {
		recurentExpensesRepo repos.RecurrentExpenseRepo
	}
)

func NewCreateRecurrentExpenseImpl(
	recurentExpensesRepo repos.RecurrentExpenseRepo,
) *CreateRecurrentExpenseImpl {
	return &CreateRecurrentExpenseImpl{
		recurentExpensesRepo: recurentExpensesRepo,
	}
}

func (c *CreateRecurrentExpenseImpl) Create(
	ctx context.Context,
	input *CreateRecurrentExpenseInput,
) (*CreateRecurrentExpenseOutput, error) {
	log.Println("Request: ", json.MustMarshall(input))
	var (
		recurrentExpense = entities.RecurrentExpense{
			Name:        input.Name,
			Amount:      input.Amount,
			Description: input.Description,
		}
	)
	if err := c.recurentExpensesRepo.Save(ctx, &recurrentExpense); err != nil {
		return nil, err
	}

	return &CreateRecurrentExpenseOutput{
		RecurrentExpense: &recurrentExpense,
	}, nil
}
