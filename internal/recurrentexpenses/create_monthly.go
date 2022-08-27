package recurrentexpenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	CreateMonthlyRecurrentExpenses interface {
		CreateMonthlyRecurrentExpenses(ctx context.Context) error
	}
	CreateMonthlyRecurrentExpensesImpl struct {
		recurrentExpensesRepo repos.RecurrentExpenseRepo
		expensesRepo          repos.ExpensesRepository
		timeGetable           dates.TimeGetable
	}
)

func NewCreateMonthlyRecurrentExpensesImpl(
	recurrentExpensesRepo repos.RecurrentExpenseRepo,
	expensesRepo repos.ExpensesRepository,
	timeGetable dates.TimeGetable,
) *CreateMonthlyRecurrentExpensesImpl {
	return &CreateMonthlyRecurrentExpensesImpl{
		recurrentExpensesRepo: recurrentExpensesRepo,
		expensesRepo:          expensesRepo,
		timeGetable:           timeGetable,
	}
}

func (c *CreateMonthlyRecurrentExpensesImpl) CreateMonthlyRecurrentExpenses(ctx context.Context) error {
	allRecurrentExpensesRegistered, err := c.recurrentExpensesRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	nextMonthDay := c.timeGetable.GetNextMonthAtFirtsDay()
	for _, recurrentExpense := range *allRecurrentExpensesRegistered {
		_, err := c.expensesRepo.FindByNameAndIsRecurrent(ctx, recurrentExpense.Name)
		if err != nil {
			_, isNotFound := err.(*repos.NotFoundError)
			if isNotFound {

				expenseToSave := entities.Expense{
					Name:        recurrentExpense.Name,
					Description: recurrentExpense.Description,
					Amount:      recurrentExpense.Amount,
					IsRecurrent: true,
					CreatedAt:   &nextMonthDay,
				}
				if err := c.expensesRepo.Save(ctx, &expenseToSave); err != nil {
					return err
				}
				continue
			}
		}

	}

	return nil
}
