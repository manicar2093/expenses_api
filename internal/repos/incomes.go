package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IncomesRepository interface {
		Save(context.Context, *entities.Income) error
	}
	IncomesRepositoryImpl struct {
		coll *mongo.Collection
	}
)

func NewIncomesRepositoryImpl(conn *mongo.Database) *IncomesRepositoryImpl {
	return &IncomesRepositoryImpl{
		coll: conn.Collection(entities.IncomeCollectionName),
	}
}

func (c *IncomesRepositoryImpl) Save(ctx context.Context, income *entities.Income) error {
	income.ID = primitive.NewObjectID()
	createdAt := dates.GetNormalizedDate()
	income.CreatedAt = &createdAt

	income.Day = uint(createdAt.Day())
	income.Month = uint(createdAt.Month())
	income.Year = uint(createdAt.Year())

	if _, err := c.coll.InsertOne(ctx, income); err != nil {
		return err
	}

	return nil
}
