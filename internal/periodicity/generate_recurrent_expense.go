package periodicity

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/periodicity/periodizer"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	ExpensePeriodicityServiceImpl struct {
		expensesRepo                        repos.ExpensesRepository
		recurrentExpensesRepo               repos.RecurrentExpenseRepo
		recurrentExpensesMonthlyCreatedRepo repos.RecurrentExpensesMonthlyCreatedRepo
		timeGetter                          dates.TimeGetable
		periodizerExpensesGenerator         periodizer.ExpensesCountByRecurrentExpensePeriodicityGenerable
		expensesCountRepo                   repos.ExpensesCountRepo
	}
)

func NewExpensePeriodicityServiceImpl(
	expensesRepo repos.ExpensesRepository,
	recurrentExpensesRepo repos.RecurrentExpenseRepo,
	recurrentExpensesMonthlyCreatedRepo repos.RecurrentExpensesMonthlyCreatedRepo,
	timeGetter dates.TimeGetable,
	periodizerExpensesGenerator periodizer.ExpensesCountByRecurrentExpensePeriodicityGenerable,
	expensesCountRepo repos.ExpensesCountRepo,
) *ExpensePeriodicityServiceImpl {
	return &ExpensePeriodicityServiceImpl{
		expensesRepo:                        expensesRepo,
		recurrentExpensesRepo:               recurrentExpensesRepo,
		recurrentExpensesMonthlyCreatedRepo: recurrentExpensesMonthlyCreatedRepo,
		timeGetter:                          timeGetter,
		periodizerExpensesGenerator:         periodizerExpensesGenerator,
		expensesCountRepo:                   expensesCountRepo,
	}
}

func (c *ExpensePeriodicityServiceImpl) GenerateRecurrentExpensesByYearAndMonth(ctx context.Context, month, year uint) (
	*entities.RecurrentExpensesMonthlyCreated,
	error,
) {
	log.Println("Getting data to:", "month:", month, "year:", year)
	savedData, err := c.recurrentExpensesMonthlyCreatedRepo.FindByCurrentMonthAndYear(ctx, month, year)
	_, notFound := err.(*repos.NotFoundError)
	switch {
	case savedData != nil:
		log.Printf("Data already saved. Returning: %+v\n", savedData)
		return savedData, nil
	case notFound:
		break
	case err != nil:
		return nil, err
	}

	log.Println("Calculating recurrent expenses...")
	recurrentExpensesRegistered, err := c.recurrentExpensesRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var (
		today          = c.timeGetter.GetCurrentTime()
		expensesCounts = make([]*entities.ExpensesCount, 0)
		response       = entities.RecurrentExpensesMonthlyCreated{
			Month: uint(today.Month()),
			Year:  uint(today.Year()),
		}
	)
	if err := c.recurrentExpensesMonthlyCreatedRepo.Save(ctx, &response); err != nil {
		return nil, err
	}
	for _, recurrentExpense := range recurrentExpensesRegistered {
		expensesCount, err := c.createExpensesCount(ctx, recurrentExpense, &response)
		if err != nil {
			return nil, err
		}
		if expensesCount == nil {
			continue
		}
		expensesCounts = append(expensesCounts, expensesCount)
	}
	response.ExpensesCount = expensesCounts

	log.Printf("Calculation result: %+v\n", response)

	return &response, nil
}

func (c *ExpensePeriodicityServiceImpl) createExpensesCount(
	ctx context.Context,
	recurrentExpense *entities.RecurrentExpense,
	recurrentExpenseMonthlyCreated *entities.RecurrentExpensesMonthlyCreated,
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
	if err := c.recurrentExpensesRepo.Update(ctx, recurrentExpense); err != nil {
		return nil, err
	}
	res := entities.ExpensesCount{
		RecurrentExpensesMonthlyCreatedID: recurrentExpenseMonthlyCreated.ID,
		RecurrentExpenseID:                recurrentExpense.ID,
		RecurrentExpense:                  recurrentExpense,
		ExpensesRelatedIDs:                insertedRes.InsertedIDs,
		TotalExpenses:                     uint(len(insertedRes.InsertedIDs)),
		TotalExpensesPaid:                 0,
	}
	if err := c.expensesCountRepo.Save(ctx, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
