package entities

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

const RecurrentExpenseCollectonName = "recurrent_expenses"

type RecurrentExpense struct {
	ID          uuid.UUID   `json:"id,omitempty" gorm:"primaryKey,->"`
	Expenses    []*Expense  `json:"expenses,omitempty"`
	Name        string      `json:"name,omitempty"`
	Amount      float64     `json:"amount,omitempty"`
	Description null.String `json:"description,omitempty"`
	CreatedAt   *time.Time  `json:"created_at,omitempty"`
	UpdatedAt   *time.Time  `json:"updated_at,omitempty"`
}

func NewRecurrentExpense(name, description string, amount float64) *RecurrentExpense {
	return &RecurrentExpense{
		Name: name,
		Description: null.NewString(
			description,
			description == "",
		),
		Amount: amount,
	}
}
