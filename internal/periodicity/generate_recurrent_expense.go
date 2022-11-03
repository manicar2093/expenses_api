package periodicity

import (
	"context"
	"log"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/converters"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/periodtypes"
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
		today            = c.timeGetter.GetCurrentTime()
		expensesToCreate = []*entities.Expense{}
		expensesCounts   = []*entities.ExpensesCount{}
	)
	for _, recurrentExpense := range recurrentExpensesRegistered {
		switch recurrentExpense.Periodicity {
		case periodtypes.Daily:
			for i := 0; i < today.Day(); i++ {
				expensesToCreate = append(
					expensesToCreate,
					converters.CreateExpenseFromRecurrentExpense(
						recurrentExpense,
						&today,
					),
				)
			}
		case periodtypes.Weekly:
			for i := 0; i < 4; i++ {
				expensesToCreate = append(
					expensesToCreate,
					converters.CreateExpenseFromRecurrentExpense(
						recurrentExpense,
						&today,
					))
			}
		case periodtypes.FourteenDays,
			periodtypes.Paydaily:
			for i := 0; i < 2; i++ {
				expensesToCreate = append(
					expensesToCreate,
					converters.CreateExpenseFromRecurrentExpense(
						recurrentExpense,
						&today,
					))
			}
		case periodtypes.Monthly,
			periodtypes.Periodicity(0):
			expensesToCreate = append(
				expensesToCreate,
				converters.CreateExpenseFromRecurrentExpense(
					recurrentExpense,
					&today,
				))
			recurrentExpense.Periodicity = periodtypes.Monthly
		default:
			continue
		}
		insertedRes, err := c.expensesRepo.SaveMany(ctx, expensesToCreate)
		if err != nil {
			return nil, err
		}
		recurrentExpense.LastCreationDate = &today
		if err := c.recurrentExpensesRepo.Update(ctx, recurrentExpense); err != nil {
			return nil, err
		}
		expensesCounts = append(expensesCounts, &entities.ExpensesCount{
			RecurrentExpenseID: recurrentExpense.ID,
			ExpensesRelated:    insertedRes.InsertedIDs,
			TotalExpenses:      uint(len(insertedRes.InsertedIDs)),
			TotalExpensesPaid:  0,
		})
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
