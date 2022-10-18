package expenses

import (
	"context"

	"github.com/manicar2093/expenses_api/pkg/converters"
)

func (c *ExpenseServiceImpl) SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error {
	id, err := converters.TurnToObjectID(input.ID)
	if err != nil {
		return err
	}
	return c.expensesRepo.UpdateIsPaidByExpenseID(ctx, id, true)
}
