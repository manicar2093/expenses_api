package expenses

import (
	"context"
	"fmt"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/json"
)

type (
	CreateExpenseInput struct {
		Name         string  `json:"name,omitempty"`
		Amount       float64 `json:"amount,omitempty"`
		Description  string  `json:"description,omitempty"`
		ForNextMonth bool    `json:"for_next_month,omitempty"`
	}
	CreateExpense interface {
		Create(context.Context, *CreateExpenseInput) (*entities.Expense, error)
	}
	CreateExpenseImpl struct {
		expensesRepo repos.ExpensesRepository
		timeGetter   dates.TimeGetable
	}
)

func NewCreateExpensesImpl(repo repos.ExpensesRepository, timeGetter dates.TimeGetable) *CreateExpenseImpl {
	return &CreateExpenseImpl{
		expensesRepo: repo,
		timeGetter:   timeGetter,
	}
}

func (c *CreateExpenseImpl) Create(ctx context.Context, expense *CreateExpenseInput) (*entities.Expense, error) {
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
