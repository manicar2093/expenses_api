package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	RecurrentExpensesMonthlyCreatedRepo interface {
		Save(
			ctx context.Context,
			recurrentExpense *entities.RecurrentExpensesMonthlyCreated,
		)
	}
	RecurrentExpensesMonthlyCreatedRepoImpl struct {
		coll *mongo.Collection
	}
)

func NewRecurrentExpensesCreatedMonthlyRepoImpl(
	conn *mongo.Database,
) *RecurrentExpensesMonthlyCreatedRepoImpl {
	return &RecurrentExpensesMonthlyCreatedRepoImpl{
		coll: conn.Collection(entities.RecurrentExpensesCreatedMonthlyCollectionName),
	}
}

func (c *RecurrentExpensesMonthlyCreatedRepoImpl) Save(
	ctx context.Context,
	recurrentExpense *entities.RecurrentExpensesMonthlyCreated,
) error {
	recurrentExpense.ID = primitive.NewObjectID()
	_, err := c.coll.InsertOne(ctx, recurrentExpense)
	if err != nil {
		return err
	}
	return nil
}
