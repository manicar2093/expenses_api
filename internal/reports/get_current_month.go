package reports

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
)

type (
	CurrentMonthDetailsOutput struct {
		TotalPaidAmount     float64            `json:"total_paid_amount,omitempty"`
		TotalUnpaidAmount   float64            `json:"total_unpaid_amount,omitempty"`
		ExpensesCount       uint               `json:"expenses_count,omitempty"`
		PaidExpensesCount   uint               `json:"paid_expenses_count,omitempty"`
		UnpaidExpensesCount uint               `json:"unpaid_expenses_count,omitempty"`
		Expenses            []entities.Expense `json:"expenses,omitempty"`
		PaidExpenses        []entities.Expense `json:"paid_expenses,omitempty"`
		UnpaidExpenses      []entities.Expense `json:"unpaid_expenses,omitempty"`
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
		paidExpenses   []entities.Expense
		unpaidExpenses []entities.Expense
		totalPaid      float64
		totalUnpaid    float64
	)

	for _, expense := range *monthExpenses {
		if expense.IsPaid {
			paidExpenses = append(paidExpenses, expense)
			totalPaid += expense.Amount
			continue
		}
		unpaidExpenses = append(unpaidExpenses, expense)
		totalUnpaid += expense.Amount
	}

	return &CurrentMonthDetailsOutput{
		TotalPaidAmount:     totalPaid,
		TotalUnpaidAmount:   totalUnpaid,
		ExpensesCount:       uint(len(*monthExpenses)),
		PaidExpensesCount:   uint(len(paidExpenses)),
		UnpaidExpensesCount: uint(len(unpaidExpenses)),
		Expenses:            *monthExpenses,
		PaidExpenses:        paidExpenses,
		UnpaidExpenses:      unpaidExpenses,
	}, nil
}
