package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	RecurrentExpenseRepo interface {
		Save(ctx context.Context, recExpense *entities.RecurrentExpense) error
	}
	RecurrentExpenseRepoImpl struct {
		coll *mongo.Collection
	}
)

func NewRecurrentExpenseRepoImpl(conn *mongo.Database) *RecurrentExpenseRepoImpl {
	return &RecurrentExpenseRepoImpl{
		coll: conn.Collection(entities.RecurrentExpenseCollectonName),
	}
}

func (c *RecurrentExpenseRepoImpl) Save(ctx context.Context, recExpense *entities.RecurrentExpense) error {
	recExpense.ID = primitive.NewObjectID()
	createdAt := dates.GetNormalizedDate()
	recExpense.CreatedAt = &createdAt
	if _, err := c.coll.InsertOne(ctx, &recExpense); err != nil {
		switch herr := err.(type) {
		case mongo.WriteException:
			if herr.WriteErrors[0].Code == 11000 {
				return &AlreadyExistsError{
					Identifier: recExpense.Name,
					Entity:     "RecurrentExpense",
				}
			}
		}
		return err
	}
	return nil
}
