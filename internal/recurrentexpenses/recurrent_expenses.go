package recurrentexpenses

import (
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

type RecurrentExpenseServiceImpl struct {
	recurrentExpensesRepo repos.RecurrentExpenseRepo
	expensesRepo          repos.ExpensesRepository
	timeGetter            dates.TimeGetable
	validator             validator.StructValidable
}

func NewCreateRecurrentExpense(
	recurentExpensesRepo repos.RecurrentExpenseRepo,
	expensesRepo repos.ExpensesRepository,
	timeGetter dates.TimeGetable,
	validator validator.StructValidable,
) *RecurrentExpenseServiceImpl {
	return &RecurrentExpenseServiceImpl{
		recurrentExpensesRepo: recurentExpensesRepo,
		expensesRepo:          expensesRepo,
		timeGetter:            timeGetter,
		validator:             validator,
	}
}

func NewGetAllRecurrentExpenseServiceImpl(
	recurentExpensesRepo repos.RecurrentExpenseRepo,
) *RecurrentExpenseServiceImpl {
	return &RecurrentExpenseServiceImpl{
		recurrentExpensesRepo: recurentExpensesRepo,
	}
}

func NewCreateMonthlyRecurrentExpensesImpl(
	recurrentExpensesRepo repos.RecurrentExpenseRepo,
	expensesRepo repos.ExpensesRepository,
	timeGettable dates.TimeGetable,
) *RecurrentExpenseServiceImpl {
	return &RecurrentExpenseServiceImpl{
		recurrentExpensesRepo: recurrentExpensesRepo,
		expensesRepo:          expensesRepo,
		timeGetter:            timeGettable,
	}
}
