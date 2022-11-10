package entityconv

import (
	"time"

	"github.com/manicar2093/expenses_api/internal/entities"
)

func CreateExpenseFromRecurrentExpense(
	recurrentExpense *entities.RecurrentExpense,
	date *time.Time,
) *entities.Expense {
	return &entities.Expense{
		RecurrentExpenseID: &recurrentExpense.ID,
		Name:               recurrentExpense.Name,
		Amount:             recurrentExpense.Amount,
		Description:        recurrentExpense.Description,
		Day:                uint(date.Day()),
		Month:              uint(date.Month()),
		Year:               uint(date.Year()),
		IsRecurrent:        true,
	}
}
