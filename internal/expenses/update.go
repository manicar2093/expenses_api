package expenses

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/pkg/json"
)

func (c *ExpenseServiceImpl) UpdateExpense(ctx context.Context, input *UpdateExpenseInput) error {
	if err := c.validator.ValidateStruct(input); err != nil {
		return err
	}
	log.Println(json.MustMarshall(input))
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
