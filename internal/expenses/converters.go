package expenses

import (
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/nullsql"
)

func createUpdateInputFromExpense(
	expense *entities.Expense,
	input *UpdateExpenseInput,
) (updateData *repos.UpdateExpenseInput) {
	if expense.IsRecurrent() {
		return &repos.UpdateExpenseInput{
			ID:     uuid.MustParse(input.ID),
			Amount: input.Amount,
		}
	}
	return &repos.UpdateExpenseInput{
		ID:          uuid.MustParse(input.ID),
		Name:        nullsql.ValidateStringSQLValid(input.Name),
		Amount:      input.Amount,
		Description: nullsql.ValidateStringSQLValid(input.Description),
	}
}
