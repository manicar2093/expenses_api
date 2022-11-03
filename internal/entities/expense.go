package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ExpenseCollectionName = "expenses"

type Expense struct {
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	RecurrentExpenses  *RecurrentExpense  `json:"recurrent_expense,omitempty" bson:"recurrent_expense,omitempty"`
	RecurrentExpenseID primitive.ObjectID `json:"recurrent_expense_id,omitempty" bson:"recurrent_expense_id,omitempty"`
	Name               string             `json:"name,omitempty"`
	Amount             float64            `json:"amount,omitempty"`
	Description        string             `json:"description,omitempty" bson:",omitempty"`
	Day                uint               `json:"day,omitempty"`
	Month              uint               `json:"month,omitempty"`
	Year               uint               `json:"year,omitempty"`
	IsRecurrent        bool               `json:"is_recurrent,omitempty" bson:"is_recurrent"`
	IsPaid             bool               `json:"is_paid,omitempty" bson:"is_paid"`
	CreatedAt          *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt          *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
