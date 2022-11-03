package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RecurrentExpensesCreatedMonthlyCollectionName = "recurrent_expenses_monthly_created"

type (
	ExpensesCount struct {
		RecurrentExpenseID primitive.ObjectID   `json:"recurrent_expense_id,omitempty" bson:"recurrent_expense_id,omitempty"`
		ExpensesRelated    []primitive.ObjectID `json:"expenses_related,omitempty" bson:"expenses_related,omitempty"`
		TotalExpenses      uint                 `json:"total_expenses,omitempty" bson:"total_expenses,omitempty"`
		TotalExpensesPaid  uint                 `json:"total_expenses_paid,omitempty" bson:"total_expenses_paid,omitempty"`
	}
	RecurrentExpensesMonthlyCreated struct {
		ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		Month         uint               `json:"month,omitempty" bson:"month,omitempty"`
		Year          uint               `json:"year,omitempty" bson:"year,omitempty"`
		ExpensesCount []*ExpensesCount   `json:"expenses_count,omitempty" bson:"expenses_count,omitempty"`
		CreatedAt     *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	}
)
