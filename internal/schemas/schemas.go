package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type ExpenseIDWithIsPaidStatus struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	IsPaid bool               `json:"is_paid,omitempty" bson:"is_paid"`
}
