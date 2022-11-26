package recurrentexpenses

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
)

type (
	CreateMonthlyRecurrentExpensesOutput struct {
		ExpensesCreated []*entities.Expense `json:"expenses_created,omitempty"`
	}
)

func (c *RecurrentExpenseServiceImpl) CreateMonthlyRecurrentExpenses(ctx context.Context) (*CreateMonthlyRecurrentExpensesOutput, error) {
	allRecurrentExpensesRegistered, err := c.recurrentExpensesRepo.FindAll(ctx)
	log.Printf("%v+", allRecurrentExpensesRegistered)
	if err != nil {
		return nil, err
	}

	nextMonthDate := c.timeGetter.GetNextMonthAtFirtsDay()
	nextMonthAsUint := uint(nextMonthDate.Month())
	var expensesCreated []*entities.Expense
	for _, recurrentExpense := range allRecurrentExpensesRegistered {
		_, err := c.expensesRepo.FindByNameAndMonthAndIsRecurrent(ctx, nextMonthAsUint, recurrentExpense.Name)
		if err != nil {
			_, isNotFound := err.(*repos.NotFoundError)
			if isNotFound {
				expenseToSave := entities.Expense{
					Amount: recurrentExpense.Amount,
					Day:    uint(nextMonthDate.Day()),
					Month:  uint(nextMonthDate.Month()),
					Year:   uint(nextMonthDate.Year()),
					RecurrentExpenseID: uuid.NullUUID{
						UUID:  recurrentExpense.ID,
						Valid: true,
					},
					CreatedAt: &nextMonthDate,
				}
				if err := c.expensesRepo.Save(ctx, &expenseToSave); err != nil {
					return nil, err
				}
				expensesCreated = append(expensesCreated, &expenseToSave)
				continue
			}
		}
	}

	return &CreateMonthlyRecurrentExpensesOutput{
		ExpensesCreated: expensesCreated,
	}, nil
}
