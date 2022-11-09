package periodicity

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/periodicity/periodizer"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/periodtypes"
)

type (
	ExpensePeriodicityServiceImpl struct {
		expensesRepo                        repos.ExpensesRepository
		recurrentExpensesRepo               repos.RecurrentExpenseRepo
		recurrentExpensesMonthlyCreatedRepo repos.RecurrentExpensesMonthlyCreatedRepo
		timeGetter                          dates.TimeGetable
		periodizerExpensesGenerator         periodizer.ExpensesCountByRecurrentExpensePeriodicityGenerable
	}
)

func NewExpensePeriodicityServiceImpl(
	expensesRepo repos.ExpensesRepository,
	recurrentExpensesRepo repos.RecurrentExpenseRepo,
	recurrentExpensesMonthlyCreatedRepo repos.RecurrentExpensesMonthlyCreatedRepo,
	timeGetter dates.TimeGetable,
	periodizerExpensesGenerator periodizer.ExpensesCountByRecurrentExpensePeriodicityGenerable,
) *ExpensePeriodicityServiceImpl {
	return &ExpensePeriodicityServiceImpl{
		expensesRepo:                        expensesRepo,
		recurrentExpensesRepo:               recurrentExpensesRepo,
		recurrentExpensesMonthlyCreatedRepo: recurrentExpensesMonthlyCreatedRepo,
		timeGetter:                          timeGetter,
		periodizerExpensesGenerator:         periodizerExpensesGenerator,
	}
}

func (c *ExpensePeriodicityServiceImpl) GenerateRecurrentExpensesByYearAndMonth(ctx context.Context, month, year uint) (
	*entities.RecurrentExpensesMonthlyCreated,
	error,
) {
	savedData, err := c.recurrentExpensesMonthlyCreatedRepo.FindByCurrentMonthAndYear(ctx, month, year)
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
		expensesCount, err := c.createExpensesCount(ctx, recurrentExpense)
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

func (c *ExpensePeriodicityServiceImpl) createExpensesCount(
	ctx context.Context,
	recurrentExpense *entities.RecurrentExpense,
) (*entities.ExpensesCount, error) {
	expensesToCreate, err := c.periodizerExpensesGenerator.GenerateExpensesCountByRecurrentExpensePeriodicity(ctx, recurrentExpense)
	switch {
	case err != nil:
		return nil, err
	case expensesToCreate == nil:
		return nil, nil

	}

	insertedRes, err := c.expensesRepo.SaveMany(
		ctx, expensesToCreate,
	)
	if err != nil {
		return nil, err
	}
	today := c.timeGetter.GetCurrentTime()
	recurrentExpense.LastCreationDate = &today
	if recurrentExpense.Periodicity == periodtypes.Empty {
		recurrentExpense.Periodicity = periodtypes.Monthly
	}
	if err := c.recurrentExpensesRepo.Update(ctx, recurrentExpense); err != nil {
		return nil, err
	}
	return &entities.ExpensesCount{
		RecurrentExpenseID: recurrentExpense.ID,
		RecurrentExpense:   recurrentExpense,
		ExpensesRelated:    insertedRes.InsertedIDs,
		TotalExpenses:      uint(len(insertedRes.InsertedIDs)),
		TotalExpensesPaid:  0,
	}, nil
}
