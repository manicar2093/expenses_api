package periodicity

import (
	"time"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/entities/entityconv"
	"github.com/manicar2093/expenses_api/pkg/periodtypes"
)

type (
	PeriodicityAction func(
		uint,
		time.Time,
		*entities.RecurrentExpense,
		TimeValidatorFunc,
	) ([]*entities.Expense, bool)

	TimeValidatorFunc func(*time.Time, *time.Time, uint) bool
)

var periodicityActionMap = map[periodtypes.Periodicity]PeriodicityAction{
	periodtypes.Periodicity(0): defaulfExpensesGenerator,
	periodtypes.Daily:          defaulfExpensesGenerator,
	periodtypes.Weekly:         defaulfExpensesGenerator,
	periodtypes.FourteenDays:   defaulfExpensesGenerator,
	periodtypes.Paydaily:       defaulfExpensesGenerator,
	periodtypes.Monthly:        defaulfExpensesGenerator,
	periodtypes.BiMonthly:      defaulfExpensesGenerator,
	periodtypes.FourMonthly:    defaulfExpensesGenerator,
	periodtypes.SixMonthly:     defaulfExpensesGenerator,
	periodtypes.Yearly:         defaulfExpensesGenerator,
}

func defaulfExpensesGenerator(
	quantity uint,
	today time.Time,
	recurrentExpense *entities.RecurrentExpense,
	timeValidator TimeValidatorFunc,
) ([]*entities.Expense, bool) {
	if timeValidation(today, recurrentExpense, timeValidator) {
		return nil, false
	}
	var expenses = iterator(quantity, today, recurrentExpense)
	if recurrentExpense.Periodicity == periodtypes.Empty {
		recurrentExpense.Periodicity = periodtypes.Monthly
	}
	return expenses, true
}

func timeValidation(
	today time.Time,
	recurrentExpense *entities.RecurrentExpense,
	timeValidator TimeValidatorFunc,
) bool {
	return timeValidator != nil && !timeValidator(
		&today,
		recurrentExpense.LastCreationDate,
		uint(recurrentExpense.Periodicity.GetTimeQuantity()),
	)
}

func iterator(
	quantity uint,
	today time.Time,
	recurrentExpense *entities.RecurrentExpense,
) []*entities.Expense {
	var expenses []*entities.Expense
	for i := 0; i < int(quantity); i++ {
		expenses = append(
			expenses,
			entityconv.CreateExpenseFromRecurrentExpense(
				recurrentExpense,
				&today,
			),
		)
	}
	return expenses
}
