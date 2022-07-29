package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecurrentExpense struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name        string             `json:"name,omitempty"`
	Amount      float64            `json:"amount,omitempty"`
	Description string             `json:"description,omitempty" bson:",omitempty"`
	CreatedAt   *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
