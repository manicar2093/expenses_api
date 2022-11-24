package recurrentexpenses

import (
	"context"
)

func (c *RecurrentExpenseServiceImpl) GetAll(ctx context.Context) (*GetAllRecurrentExpensesOutput, error) {
	recurrentExpenses, err := c.recurrentExpensesRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return &GetAllRecurrentExpensesOutput{
		RecurrentExpenses:       recurrentExpenses,
		RecurrenteExpensesCount: uint(len(recurrentExpenses)),
	}, nil
}
