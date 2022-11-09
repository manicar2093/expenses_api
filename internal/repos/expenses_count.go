package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ExpensesCountRepo interface {
		Save(ctx context.Context, expenseCount *entities.ExpensesCount) error
	}

	ExpensesCountMongoRepo struct {
		coll *mongo.Collection
	}
)

func NewExpensesCountMongoRepo(conn *mongo.Database) *ExpensesCountMongoRepo {
	return &ExpensesCountMongoRepo{coll: conn.Collection(entities.ExpensesCountCollName)}
}

func (c *ExpensesCountMongoRepo) Save(ctx context.Context, expenseCount *entities.ExpensesCount) error {
	expenseCount.ID = primitive.NewObjectID()
	_, err := c.coll.InsertOne(ctx, &expenseCount)
	if err != nil {
		return err
	}
	return nil
}
