package reports

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/converters"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	CurrentMonthDetailsOutput struct {
		TotalPaidAmount     float64                  `json:"total_paid_amount"`
		TotalUnpaidAmount   float64                  `json:"total_unpaid_amount"`
		ExpensesCount       uint                     `json:"expenses_count"`
		PaidExpensesCount   uint                     `json:"paid_expenses_count"`
		UnpaidExpensesCount uint                     `json:"unpaid_expenses_count"`
		Expenses            []*mongoentities.Expense `json:"expenses"`
		PaidExpenses        []*mongoentities.Expense `json:"paid_expenses"`
		UnpaidExpenses      []*mongoentities.Expense `json:"unpaid_expenses"`
	}
	CurrentMonthDetailsGettable interface {
		GetExpenses(ctx context.Context) (*CurrentMonthDetailsOutput, error)
	}
	CurrentMonthDetails struct {
		repo       repos.ExpensesRepository
		timeGetter dates.TimeGetable
	}
)

func NewCurrentMonthDetailsImpl(repo repos.ExpensesRepository, timeGetter dates.TimeGetable) *CurrentMonthDetails {
	return &CurrentMonthDetails{repo: repo, timeGetter: timeGetter}
}

func (c *CurrentMonthDetails) GetExpenses(ctx context.Context) (*CurrentMonthDetailsOutput, error) {
	monthExpenses, err := c.repo.GetExpensesByMonth(ctx, c.timeGetter.GetCurrentTime().Month())
	if err != nil {
		return nil, err
	}

	var (
		paidExpenses   = []*mongoentities.Expense{}
		unpaidExpenses = []*mongoentities.Expense{}
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

	return &CurrentMonthDetailsOutput{
		TotalPaidAmount:     converters.Round(totalPaid),
		TotalUnpaidAmount:   converters.Round(totalUnpaid),
		ExpensesCount:       uint(len(monthExpenses)),
		PaidExpensesCount:   uint(len(paidExpenses)),
		UnpaidExpensesCount: uint(len(unpaidExpenses)),
		Expenses:            monthExpenses,
		PaidExpenses:        paidExpenses,
		UnpaidExpenses:      unpaidExpenses,
	}, nil
}
