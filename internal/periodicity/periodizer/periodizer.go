package periodizer

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	PeriodicityService struct {
		expensesRepo         repos.ExpensesRepository
		recurrenExpensesRepo repos.RecurrentExpenseRepo
		timeGetter           dates.TimeGetable
	}
)

func NewPeriodicityService(
	expensesRepo repos.ExpensesRepository,
	recurrenExpensesRepo repos.RecurrentExpenseRepo,
	timeGetter dates.TimeGetable,
) *PeriodicityService {
	return &PeriodicityService{
		expensesRepo:         expensesRepo,
		recurrenExpensesRepo: recurrenExpensesRepo,
		timeGetter:           timeGetter,
	}
}

func (c *PeriodicityService) GenerateExpensesCountByRecurrentExpensePeriodicity(
	ctx context.Context,
	recurrentExpense *entities.RecurrentExpense,
) ([]*entities.Expense, error) {
	expensesToCreate, hasExpenses := c.categorizeByPeriodicity(ctx, recurrentExpense)
	if !hasExpenses {
		return nil, nil
	}
	return expensesToCreate, nil
}

func (c *PeriodicityService) categorizeByPeriodicity(
	ctx context.Context,
	recurrentExpense *entities.RecurrentExpense,
) ([]*entities.Expense, bool) {
	var (
		today       = c.timeGetter.GetCurrentTime()
		periodicity = recurrentExpense.Periodicity
	)
	action, ok := periodicityActionMap[periodicity]
	if !ok {
		log.Printf(
			"unhandled periodicity for '%s' recurrent expense: '%v'",
			recurrentExpense.Name,
			recurrentExpense.Periodicity,
		)
		return nil, false
	}
	if recurrentExpense.LastCreationDate == nil {
		return action(periodicity.GetExpensesQuantity(today), today, recurrentExpense, nil)
	}
	return action(periodicity.GetExpensesQuantity(today), today, recurrentExpense, periodicity.GetTimeValidator())
}
