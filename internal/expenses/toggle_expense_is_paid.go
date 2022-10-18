package expenses

import (
	"context"

	"github.com/manicar2093/expenses_api/pkg/json"
)

func (c *ExpenseServiceImpl) ToggleIsPaid(
	ctx context.Context,
	input *ToggleExpenseIsPaidInput,
) (*ToggleExpenseIsPaidOutput, error) {
	log.Println(json.MustMarshall(input))
	var expenseID = input.ID
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
