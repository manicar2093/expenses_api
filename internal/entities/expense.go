package entities

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

const ExpensesTableName = "expenses"

type Expense struct {
	ID                 uuid.UUID         `json:"id,omitempty" gorm:"primaryKey,->"`
	RecurrentExpense   *RecurrentExpense `json:"recurrent_expense,omitempty"`
	RecurrentExpenseID uuid.NullUUID     `json:"recurrent_expense_id,omitempty"`
	Name               string            `json:"name,omitempty"`
	Amount             float64           `json:"amount,omitempty"`
	Description        null.String       `json:"description,omitempty"`
	Day                uint              `json:"day,omitempty"`
	Month              uint              `json:"month,omitempty"`
	Year               uint              `json:"year,omitempty"`
	IsPaid             bool              `json:"is_paid,omitempty"`
	CreatedAt          *time.Time        `json:"created_at,omitempty"`
	UpdatedAt          *time.Time        `json:"updated_at,omitempty"`
}

type ExpenseIDWithIsPaidStatus struct {
	ID     uuid.UUID `json:"id,omitempty"`
	IsPaid bool      `json:"is_paid"`
}

func NewExpense(name, description string, amount float64, recurrentExpense *RecurrentExpense, createAt *time.Time) *Expense {
	var recurrentExpenseID uuid.NullUUID
	if recurrentExpense != nil {
		recurrentExpenseID = uuid.NullUUID{
			UUID:  recurrentExpense.ID,
			Valid: true,
		}
	}
	return &Expense{
		Name:   name,
		Amount: amount,
		Description: null.NewString(
			description,
			description != "",
		),
		RecurrentExpenseID: recurrentExpenseID,
		CreatedAt:          createAt,
	}
}
