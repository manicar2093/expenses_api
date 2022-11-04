package periodicity

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	ExpensePeriodicityServiceImpl struct {
		expensesRepo                        repos.ExpensesRepository
		recurrentExpensesRepo               repos.RecurrentExpenseRepo
		recurrentExpensesMonthlyCreatedRepo repos.RecurrentExpensesMonthlyCreatedRepo
		timeGetter                          dates.TimeGetable
	}
)

func NewExpensePeriodicityServiceImpl(
	expensesRepo repos.ExpensesRepository,
	recurrentExpensesRepo repos.RecurrentExpenseRepo,
	recurrentExpensesMonthlyCreatedRepo repos.RecurrentExpensesMonthlyCreatedRepo,
	timeGetter dates.TimeGetable,
) *ExpensePeriodicityServiceImpl {
	return &ExpensePeriodicityServiceImpl{
		expensesRepo:                        expensesRepo,
		recurrentExpensesRepo:               recurrentExpensesRepo,
		recurrentExpensesMonthlyCreatedRepo: recurrentExpensesMonthlyCreatedRepo,
		timeGetter:                          timeGetter,
	}
}

func (c *ExpensePeriodicityServiceImpl) GenerateRecurrentExpensesByYearAndMonth(ctx context.Context, month, year uint) (
	*entities.RecurrentExpensesMonthlyCreated,
	error,
) {
	savedData, err := c.recurrentExpensesMonthlyCreatedRepo.FindByMonthAndYear(ctx, month, year)
	_, notFound := err.(*repos.NotFoundError)
	switch {
	case savedData != nil:
		return savedData, nil
	case notFound:
		break
	case err != nil:
		return nil, err
	}

	recurrentExpensesRegistered, err := c.recurrentExpensesRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var (
		today          = c.timeGetter.GetCurrentTime()
		expensesCounts = []*entities.ExpensesCount{}
	)
	for _, recurrentExpense := range recurrentExpensesRegistered {
		expensesCount, err := c.GenerateExpensesByPeriodicity(ctx, recurrentExpense)
		if err != nil {
			return nil, err
		}
		if expensesCount == nil {
			continue
		}
		expensesCounts = append(expensesCounts, expensesCount)
	}

	response := entities.RecurrentExpensesMonthlyCreated{
		Month:         uint(today.Month()),
		Year:          uint(today.Year()),
		ExpensesCount: expensesCounts,
	}
	if err := c.recurrentExpensesMonthlyCreatedRepo.Save(ctx, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *ExpensePeriodicityServiceImpl) GenerateExpensesByPeriodicity(
	ctx context.Context,
	recurrentExpense *entities.RecurrentExpense,
) (*entities.ExpensesCount, error) {
	var today = c.timeGetter.GetCurrentTime()
	expensesToCreate, hasExpenses := c.CategorizeByPeriodicity(ctx, recurrentExpense)
	if !hasExpenses {
		return nil, nil
	}
	insertedRes, err := c.expensesRepo.SaveMany(
		ctx, expensesToCreate,
	)
	if err != nil {
		return nil, err
	}
	recurrentExpense.LastCreationDate = &today
	if err := c.recurrentExpensesRepo.Update(ctx, recurrentExpense); err != nil {
		return nil, err
	}
	return &entities.ExpensesCount{
		RecurrentExpenseID: recurrentExpense.ID,
		ExpensesRelated:    insertedRes.InsertedIDs,
		TotalExpenses:      uint(len(insertedRes.InsertedIDs)),
		TotalExpensesPaid:  0,
	}, nil
}

func (c *ExpensePeriodicityServiceImpl) CategorizeByPeriodicity(
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

	return action(periodicity.GetExpensesQuantity(today), today, recurrentExpense, periodicity.GetTimeValidator())
}
