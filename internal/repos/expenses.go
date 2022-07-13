package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ExpensesRepository interface {
		Save(ctx context.Context, expense *entities.Expense) error
	}
	ExpensesRepositoryImpl struct {
		coll *mongo.Collection
	}
)

func NewExpensesRepositoryImpl(coll *mongo.Database) *ExpensesRepositoryImpl {
	return &ExpensesRepositoryImpl{coll: coll.Collection("expenses")}
}

func (c *ExpensesRepositoryImpl) Save(ctx context.Context, expense *entities.Expense) error {
	expense.ID = primitive.NewObjectID()
	createdAt := dates.GetNormalizedDate()
	expense.CreatedAt = &createdAt
	if _, err := c.coll.InsertOne(ctx, expense); err != nil {
		return err
	}

	return nil
}
