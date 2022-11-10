package mongorepos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	RecurrentExpenseMongoRepo struct {
		coll *mongo.Collection
	}
)

func NewRecurrentExpenseMongoRepo(conn *mongo.Database) *RecurrentExpenseMongoRepo {
	return &RecurrentExpenseMongoRepo{
		coll: conn.Collection(mongoentities.RecurrentExpenseCollectonName),
	}
}

func (c *RecurrentExpenseMongoRepo) Save(ctx context.Context, recExpense *mongoentities.RecurrentExpense) error {
	recExpense.ID = primitive.NewObjectID()
	createdAt := dates.GetNormalizedDate()
	recExpense.CreatedAt = &createdAt
	if _, err := c.coll.InsertOne(ctx, &recExpense); err != nil {
		switch herr := err.(type) { //nolint: gocritic
		case mongo.WriteException:
			if herr.WriteErrors[0].Code == 11000 {
				return &repos.AlreadyExistsError{
					Identifier: recExpense.Name,
					Entity:     "RecurrentExpense",
				}
			}
		}
		return err
	}
	return nil
}

func (c *RecurrentExpenseMongoRepo) FindByName(ctx context.Context, name string) (*mongoentities.RecurrentExpense, error) {
	var result mongoentities.RecurrentExpense
	if err := c.coll.FindOne(ctx, bson.D{{Key: "name", Value: name}}).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *RecurrentExpenseMongoRepo) FindAll(ctx context.Context) (*[]mongoentities.RecurrentExpense, error) {
	cursor, err := c.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var result []mongoentities.RecurrentExpense
	for cursor.Next(ctx) {
		var temp mongoentities.RecurrentExpense
		if err := cursor.Decode(&temp); err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return &result, nil
}
