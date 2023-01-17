package expenses

import (
	"context"

	"github.com/google/uuid"
)

func (c *ExpenseServiceImpl) SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error {
	return c.expensesRepo.UpdateIsPaidByExpenseID(ctx, uuid.MustParse(input.ID), true)
}
