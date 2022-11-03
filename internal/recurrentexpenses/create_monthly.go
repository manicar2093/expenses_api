package recurrentexpenses

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	CreateMonthlyRecurrentExpenses interface {
		CreateMonthlyRecurrentExpenses(ctx context.Context) (*CreateMonthlyRecurrentExpensesOutput, error)
	}

	CreateMonthlyRecurrentExpensesOutput struct {
		ExpensesCreated []entities.Expense `json:"expenses_created,omitempty"`
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

func (c *CreateMonthlyRecurrentExpensesImpl) CreateMonthlyRecurrentExpenses(ctx context.Context) (*CreateMonthlyRecurrentExpensesOutput, error) {
	allRecurrentExpensesRegistered, err := c.recurrentExpensesRepo.FindAll(ctx)
	log.Printf("%v+", allRecurrentExpensesRegistered)
	if err != nil {
		return nil, err
	}

	nextMonthDate := c.timeGetable.GetNextMonthAtFirtsDay()
	nextMonthAsUint := uint(nextMonthDate.Month())
	var expensesCreated []entities.Expense
	for _, recurrentExpense := range allRecurrentExpensesRegistered {
		_, err := c.expensesRepo.FindByNameAndMonthAndIsRecurrent(ctx, nextMonthAsUint, recurrentExpense.Name)
		if err != nil {
			_, isNotFound := err.(*repos.NotFoundError)
			if isNotFound {
				expenseToSave := entities.Expense{
					Name:        recurrentExpense.Name,
					Description: recurrentExpense.Description,
					Amount:      recurrentExpense.Amount,
					CreatedAt:   &nextMonthDate,
				}
				if err := c.expensesRepo.SaveAsRecurrent(ctx, &expenseToSave); err != nil {
					return nil, err
				}
				expensesCreated = append(expensesCreated, expenseToSave)
				continue
			}
		}
	}

	return &CreateMonthlyRecurrentExpensesOutput{
		ExpensesCreated: expensesCreated,
	}, nil
}
