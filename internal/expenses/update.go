package expenses

import (
	"context"

	"github.com/google/uuid"
)

func (c *ExpenseServiceImpl) UpdateExpense(ctx context.Context, input *UpdateExpenseInput) error {
	var expenseID = uuid.MustParse(input.ID)
	stored, err := c.expensesRepo.FindByID(ctx, expenseID)
	if err != nil {
		return err
	}
	if err := c.expensesRepo.Update(ctx, createUpdateInputFromExpense(stored, input)); err != nil {
		return err
	}

	return nil
}
