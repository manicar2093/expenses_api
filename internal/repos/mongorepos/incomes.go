package mongorepos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IncomesMongoRepo struct {
		coll *mongo.Collection
	}
)

func NewIncomesMongoRepo(conn *mongo.Database) *IncomesMongoRepo {
	return &IncomesMongoRepo{
		coll: conn.Collection(mongoentities.IncomeCollectionName),
	}
}

func (c *IncomesMongoRepo) Save(ctx context.Context, income *mongoentities.Income) error {
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
