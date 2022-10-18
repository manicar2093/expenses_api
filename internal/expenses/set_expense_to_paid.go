package expenses

import (
	"context"
)

func (c *ExpenseServiceImpl) SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error {
	return c.expensesRepo.UpdateIsPaidByExpenseID(ctx, input.ID, true)
}
