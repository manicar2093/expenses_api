package entities

import (
	"time"

	"github.com/manicar2093/expenses_api/pkg/periodtypes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RecurrentExpensesCollectonName = "recurrent_expenses"

type RecurrentExpense struct {
	ID               primitive.ObjectID      `json:"id,omitempty" bson:"_id"`
	Expenses         []*Expense              `json:"expenses,omitempty" bson:",omitempty"`
	Name             string                  `json:"name,omitempty"`
	Amount           float64                 `json:"amount,omitempty"`
	Description      string                  `json:"description,omitempty" bson:",omitempty"`
	Periodicity      periodtypes.Periodicity `json:"periodicity,omitempty"`
	LastCreationDate *time.Time              `json:"last_creation_date,omitempty" bson:"last_creation_date,omitempty"`
	CreatedAt        *time.Time              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt        *time.Time              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
