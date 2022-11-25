package recurrentexpenses

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/json"
	"github.com/manicar2093/expenses_api/pkg/nullsql"
)

func (c *RecurrentExpenseServiceImpl) CreateRecurrentExpense(
	ctx context.Context,
	input *CreateRecurrentExpenseInput,
) (*CreateRecurrentExpenseOutput, error) {
	if err := c.validator.ValidateStruct(input); err != nil {
		return nil, err
	}
	log.Println("Request: ", json.MustMarshall(input))
	var (
		nextMontTime     = c.timeGetter.GetNextMonthAtFirtsDay()
		recurrentExpense = entities.RecurrentExpense{
			Name:        input.Name,
			Description: nullsql.ValidateStringSQLValid(input.Description),
			Amount:      input.Amount,
		}
		expense = entities.Expense{
			Amount:    input.Amount,
			CreatedAt: &nextMontTime,
		}
	)
	if err := c.recurrentExpensesRepo.Save(ctx, &recurrentExpense); err != nil {
		return nil, err
	}
	expense.RecurrentExpenseID = uuid.NullUUID{
		UUID:  recurrentExpense.ID,
		Valid: true,
	}
	if err := c.expensesRepo.Save(ctx, &expense); err != nil {
		return nil, err
	}

	return &CreateRecurrentExpenseOutput{
		RecurrentExpense: &recurrentExpense,
		NextMonthExpense: &expense,
	}, nil
}
