package entities

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

const (
	ExpensesTableName  = "expenses"
	ExpensesEntityName = "expense"
)

type Expense struct {
	ID                 uuid.UUID         `json:"id,omitempty" gorm:"primaryKey,->"`
	RecurrentExpense   *RecurrentExpense `json:"recurrent_expense,omitempty"`
	RecurrentExpenseID uuid.NullUUID     `json:"recurrent_expense_id,omitempty"`
	User               *User             `json:"user,omitempty"`
	UserID             uuid.UUID         `json:"user_id,omitempty"`
	Name               null.String       `json:"name,omitempty"`
	Amount             float64           `json:"amount,omitempty"`
	Description        null.String       `json:"description,omitempty"`
	Day                uint              `json:"day,omitempty"`
	Month              uint              `json:"month,omitempty"`
	Year               uint              `json:"year,omitempty"`
	IsPaid             bool              `json:"is_paid"`
	CreatedAt          *time.Time        `json:"created_at,omitempty"`
	UpdatedAt          *time.Time        `json:"updated_at,omitempty"  gorm:"autoUpdateTime:false"`
}

type ExpenseIDWithIsPaidStatus struct {
	ID     uuid.UUID `json:"id,omitempty"`
	IsPaid bool      `json:"is_paid"`
}

func (c *Expense) IsRecurrent() bool {
	return c.RecurrentExpenseID.Valid
}
