package expenses

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/nullsql"
	"gopkg.in/guregu/null.v4"
)

func (c *ExpenseServiceImpl) CreateExpense(ctx context.Context, expense *CreateExpenseInput) (*entities.Expense, error) {
	newExpense := c.expenseFromCreateExpenseInput(expense)
	if err := c.expensesRepo.Save(ctx, newExpense); err != nil {
		return nil, err
	}
	return newExpense, nil
}

func (c *ExpenseServiceImpl) expenseFromCreateExpenseInput(holder *CreateExpenseInput) *entities.Expense {
	var (
		today      = c.timaGettable.GetCurrentTime()
		newExpense = &entities.Expense{
			Name:        nullsql.ValidateStringSQLValid(holder.Name),
			Amount:      holder.Amount,
			Description: nullsql.ValidateStringSQLValid(holder.Description),
			IsPaid:      true,
			CreatedAt:   &today,
			UserID:      uuid.MustParse(holder.UserID),
		}
	)
	if holder.ForNextMonth {
		nextMonthTime := c.timaGettable.GetNextMonthAtFirtsDay()
		newExpense.Description = null.StringFrom(
			fmt.Sprintf(
				"%s\n\nFecha de registro: %s",
				newExpense.Description.ValueOrZero(),
				today.Format("02/01/2006"),
			),
		)
		newExpense.CreatedAt = &nextMonthTime
	}
	newExpense.Day = uint(newExpense.CreatedAt.Day())
	newExpense.Month = uint(newExpense.CreatedAt.Month())
	newExpense.Year = uint(newExpense.CreatedAt.Year())

	return newExpense
}
