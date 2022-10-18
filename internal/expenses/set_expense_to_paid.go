package expenses

import (
	"context"
)

func (c *ExpenseServiceImpl) SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error {
	log.Println("DEPRECATED!. Use ExpenseToPaidTogglable instead")
	return c.expensesRepo.UpdateIsPaidByExpenseID(ctx, input.ID, true)
}
