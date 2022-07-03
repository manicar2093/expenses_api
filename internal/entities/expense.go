package entities

import "time"

type Expense struct {
	ID          uint      `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (c Expense) Table() string {
	return "Expenses"
}
