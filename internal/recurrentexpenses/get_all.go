package recurrentexpenses

import (
	"context"

	"github.com/google/uuid"
)

func (c *RecurrentExpenseServiceImpl) GetAll(ctx context.Context, userID uuid.UUID) (*GetAllRecurrentExpensesOutput, error) {
	recurrentExpenses, err := c.recurrentExpensesRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetAllRecurrentExpensesOutput{
		RecurrentExpenses:       recurrentExpenses,
		RecurrenteExpensesCount: uint(len(recurrentExpenses)),
	}, nil
}
