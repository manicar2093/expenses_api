package expenses

import (
	"context"
	"fmt"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/json"
)

func (c *ExpenseServiceImpl) Create(ctx context.Context, expense *CreateExpenseInput) (*entities.Expense, error) {
	log.Println(json.MustMarshall(expense))
	newExpense := &entities.Expense{
		Name:        expense.Name,
		Amount:      expense.Amount,
		Description: expense.Description,
		IsPaid:      true,
	}
	if expense.ForNextMonth {
		nextMonthTime := c.timeGetter.GetNextMonthAtFirtsDay()
		today := c.timeGetter.GetCurrentTime()
		newExpense.Description = fmt.Sprintf(
			"%s\n\nFecha de registro: %s",
			newExpense.Description,
			today.Format("02/01/2006"),
		)
		newExpense.CreatedAt = &nextMonthTime
	}
	if err := c.expensesRepo.Save(ctx, newExpense); err != nil {
		return nil, err
	}
	return newExpense, nil
}
