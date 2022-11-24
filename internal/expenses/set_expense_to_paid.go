package expenses

import (
	"context"

	"github.com/google/uuid"
)

func (c *ExpenseServiceImpl) SetToPaid(ctx context.Context, input *SetExpenseToPaidInput) error {
	log.Println("DEPRECATED!. Use ExpenseToPaidTogglable instead")
	if err := c.validator.ValidateStruct(input); err != nil {
		return err
	}
	return c.expensesRepo.UpdateIsPaidByExpenseID(ctx, uuid.MustParse(input.ID), true)
}
