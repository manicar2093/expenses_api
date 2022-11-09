package reports

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/periodicity"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/converters"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	CurrentMonthDetailsOutput struct {
		TotalPaidAmount          float64                                   `json:"total_paid_amount"`
		TotalUnpaidAmount        float64                                   `json:"total_unpaid_amount"`
		ExpensesCount            uint                                      `json:"expenses_count"`
		PaidExpensesCount        uint                                      `json:"paid_expenses_count"`
		UnpaidExpensesCount      uint                                      `json:"unpaid_expenses_count"`
		Expenses                 []*entities.Expense                       `json:"expenses"`
		PaidExpenses             []*entities.Expense                       `json:"paid_expenses"`
		UnpaidExpenses           []*entities.Expense                       `json:"unpaid_expenses"`
		RecurrentExpensesDetails *entities.RecurrentExpensesMonthlyCreated `json:"recurrent_expenses_details"`
	}
	CurrentMonthDetailsGettable interface {
		GetCurrentMonthDetails(ctx context.Context) (*CurrentMonthDetailsOutput, error)
	}
	CurrentMonthDetails struct {
		repo                           repos.ExpensesRepository
		periodicityMonthAndYearCreator periodicity.ExpensePeriodicityCreateable
		timeGetter                     dates.TimeGetable
	}
)

func NewCurrentMonthDetailsImpl(
	repo repos.ExpensesRepository,
	timeGetter dates.TimeGetable,
	periodicityMonthAndYearCreator periodicity.ExpensePeriodicityCreateable,
) *CurrentMonthDetails {
	return &CurrentMonthDetails{repo: repo, timeGetter: timeGetter, periodicityMonthAndYearCreator: periodicityMonthAndYearCreator}
}

func (c *CurrentMonthDetails) GetCurrentMonthDetails(ctx context.Context) (*CurrentMonthDetailsOutput, error) {
	var (
		month = c.timeGetter.GetCurrentTime().Month()
		year  = c.timeGetter.GetCurrentTime().Year()
	)
	monthExpenses, err := c.repo.GetExpensesByMonth(ctx, month)
	if err != nil {
		return nil, err
	}

	var (
		paidExpenses   = []*entities.Expense{}
		unpaidExpenses = []*entities.Expense{}
		totalPaid      float64
		totalUnpaid    float64
	)

	for _, expense := range monthExpenses {
		if expense.IsPaid {
			paidExpenses = append(paidExpenses, expense)
			totalPaid += expense.Amount
			continue
		}
		unpaidExpenses = append(unpaidExpenses, expense)
		totalUnpaid += expense.Amount
	}

	recurrentExpensesDetails, err := c.periodicityMonthAndYearCreator.GenerateRecurrentExpensesByYearAndMonth(ctx, uint(month), uint(year))
	if err != nil {
		return nil, err
	}

	return &CurrentMonthDetailsOutput{
		TotalPaidAmount:          converters.Round(totalPaid),
		TotalUnpaidAmount:        converters.Round(totalUnpaid),
		ExpensesCount:            uint(len(monthExpenses)),
		PaidExpensesCount:        uint(len(paidExpenses)),
		UnpaidExpensesCount:      uint(len(unpaidExpenses)),
		Expenses:                 monthExpenses,
		PaidExpenses:             paidExpenses,
		UnpaidExpenses:           unpaidExpenses,
		RecurrentExpensesDetails: recurrentExpensesDetails,
	}, nil
}
