package expenses

import (
	"context"

	"github.com/google/uuid"
)

func (c *ExpenseServiceImpl) ToggleIsPaid(
	ctx context.Context,
	input *ToggleExpenseIsPaidInput,
) (*ToggleExpenseIsPaidOutput, error) {
	return c.toggle(ctx, uuid.MustParse(input.ID))
}

func (c *ExpenseServiceImpl) toggle(
	ctx context.Context,
	expenseID uuid.UUID,
) (*ToggleExpenseIsPaidOutput, error) {
	expenseStatus, err := c.expensesRepo.GetExpenseStatusByID(ctx, expenseID)
	if err != nil {
		return nil, err
	}

	var newExpenseIsPaidStatus = !expenseStatus.IsPaid

	if err := c.expensesRepo.UpdateIsPaidByExpenseID(ctx, expenseID, newExpenseIsPaidStatus); err != nil {
		return nil, err
	}

	return &ToggleExpenseIsPaidOutput{
		ID:                  expenseID,
		CurrentIsPaidStatus: newExpenseIsPaidStatus,
	}, nil
}
