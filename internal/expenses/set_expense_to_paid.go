package expenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/converters"
)

type (
	SetExpenseToPaid interface {
		SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error
	}
	SetExpenseToPaidInput struct {
		ID string `json:"id,omitempty"`
	}
	SetExpenseToPaidImpl struct {
		expensesRepo repos.ExpensesRepository
	}
)

func NewSetExpenseToPaidImpl(expensesRepo repos.ExpensesRepository) *SetExpenseToPaidImpl {
	return &SetExpenseToPaidImpl{
		expensesRepo: expensesRepo,
	}
}

func (c *SetExpenseToPaidImpl) SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error {
	id, err := converters.TurnToObjectID(input.ID)
	if err != nil {
		return err
	}
	return c.expensesRepo.UpdateIsPaidByExpenseID(ctx, id, true)
}
