package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	RecurrentExpenseRepo interface {
		Save(ctx context.Context, recExpense *entities.RecurrentExpense) error
		FindByName(ctx context.Context, name string) (*entities.RecurrentExpense, error)
		FindAll(ctx context.Context) (*[]entities.RecurrentExpense, error)
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
		switch herr := err.(type) { //nolint: gocritic
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

func (c *RecurrentExpenseRepoImpl) FindByName(ctx context.Context, name string) (*entities.RecurrentExpense, error) {
	var result entities.RecurrentExpense
	if err := c.coll.FindOne(ctx, bson.D{{Key: "name", Value: name}}).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *RecurrentExpenseRepoImpl) FindAll(ctx context.Context) (*[]entities.RecurrentExpense, error) {
	cursor, err := c.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var result []entities.RecurrentExpense
	for cursor.Next(ctx) {
		var temp entities.RecurrentExpense
		if err := cursor.Decode(&temp); err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return &result, nil
}
