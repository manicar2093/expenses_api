package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

const ExpensesCountCollName = "expenses_count"

type ExpensesCount struct {
	ID                                primitive.ObjectID               `json:"id" bson:"_id"`
	RecurrentExpenseID                primitive.ObjectID               `json:"recurrent_expense_id,omitempty" bson:"recurrent_expense_id,omitempty"`
	RecurrentExpense                  *RecurrentExpense                `json:"recurrent_expense,omitempty" bson:"recurrent_expense,omitempty"`
	RecurrentExpensesMonthlyCreatedID primitive.ObjectID               `json:"recurrent_expenses_monthly_created_id" bson:"recurrent_expenses_monthly_created_id"`
	RecurrentExpensesMonthlyCreated   *RecurrentExpensesMonthlyCreated `json:"recurrent_expenses_monthly_created,omitempty" bson:"-"`
	ExpensesRelatedIDs                []primitive.ObjectID             `json:"expenses_related_ids,omitempty" bson:"expenses_related,omitempty"`
	ExpensesRelated                   []*Expense                       `json:"expenses_related,omitempty" bson:"-"`
	TotalExpenses                     uint                             `json:"total_expenses" bson:"total_expenses"`
	TotalExpensesPaid                 uint                             `json:"total_expenses_paid" bson:"total_expenses_paid"`
}
