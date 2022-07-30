package recurrentexpenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/json"
)

type (
	CreateRecurrentExpense interface {
		Create(ctx context.Context, input *CreateRecurrentExpenseInput) (*CreateRecurrentExpenseOutput, error)
	}
	CreateRecurrentExpenseInput struct {
		Name        string  `json:"name,omitempty"`
		Amount      float64 `json:"amount,omitempty"`
		Description string  `json:"description,omitempty"`
	}
	CreateRecurrentExpenseOutput struct {
		RecurrentExpense *entities.RecurrentExpense `json:"recurrent_expense,omitempty"`
		NextMonthExpense *entities.Expense          `json:"next_month_expense,omitempty"`
	}
	CreateRecurrentExpenseImpl struct {
		recurentExpensesRepo repos.RecurrentExpenseRepo
		expensesRepo         repos.ExpensesRepository
		timeGetter           dates.TimeGetable
	}
)

func NewCreateRecurrentExpenseImpl(
	recurentExpensesRepo repos.RecurrentExpenseRepo,
	expensesRepo repos.ExpensesRepository,
	timeGetter dates.TimeGetable,
) *CreateRecurrentExpenseImpl {
	return &CreateRecurrentExpenseImpl{
		recurentExpensesRepo: recurentExpensesRepo,
		expensesRepo:         expensesRepo,
		timeGetter:           timeGetter,
	}
}

func (c *CreateRecurrentExpenseImpl) Create(
	ctx context.Context,
	input *CreateRecurrentExpenseInput,
) (*CreateRecurrentExpenseOutput, error) {
	log.Println("Request: ", json.MustMarshall(input))
	var (
		nextMontTime     = c.timeGetter.GetNextMonthAtFirtsDay()
		recurrentExpense = entities.RecurrentExpense{
			Name:        input.Name,
			Amount:      input.Amount,
			Description: input.Description,
		}
		expense = entities.Expense{
			Name:        input.Name,
			Amount:      input.Amount,
			Description: input.Description,
			IsRecurrent: true,
			CreatedAt:   &nextMontTime,
		}
	)
	if err := c.recurentExpensesRepo.Save(ctx, &recurrentExpense); err != nil {
		return nil, err
	}
	if err := c.expensesRepo.Save(ctx, &expense); err != nil {
		return nil, err
	}

	return &CreateRecurrentExpenseOutput{
		RecurrentExpense: &recurrentExpense,
		NextMonthExpense: &expense,
	}, nil
}
