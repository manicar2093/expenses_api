package mongorepos

import (
	"context"
	"errors"
	"time"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/internal/schemas"
	"github.com/manicar2093/expenses_api/pkg/converters"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	ExpensesMongoRepo struct {
		coll *mongo.Collection
	}
)

func NewExpensesMongoRepo(coll *mongo.Database) *ExpensesMongoRepo {
	return &ExpensesMongoRepo{coll: coll.Collection(mongoentities.ExpenseCollectionName)}
}

func (c *ExpensesMongoRepo) Save(ctx context.Context, expense *mongoentities.Expense) error {
	expense.ID = primitive.NewObjectID()
	if expense.CreatedAt == nil {
		createdAt := dates.GetNormalizedDate()
		expense.CreatedAt = &createdAt
	}

	expense.Day = uint(expense.CreatedAt.Day())
	expense.Month = uint(expense.CreatedAt.Month())
	expense.Year = uint(expense.CreatedAt.Year())

	if _, err := c.coll.InsertOne(ctx, expense); err != nil {
		return err
	}

	return nil
}

func (c *ExpensesMongoRepo) GetExpensesByMonth(ctx context.Context, month time.Month) ([]*mongoentities.Expense, error) {
	cursor, err := c.coll.Find(ctx, bson.D{
		{Key: "month", Value: month},
	})
	if err != nil {
		return nil, err
	}

	response := make([]*mongoentities.Expense, 0)

	for cursor.Next(ctx) {
		var entityTemp mongoentities.Expense
		if err := cursor.Decode(&entityTemp); err != nil {
			return nil, err
		}
		response = append(response, &entityTemp)
	}

	return response, nil
}

func (c *ExpensesMongoRepo) UpdateIsPaidByExpenseID(ctx context.Context, expenseID interface{}, status bool) error {
	expenseObjectID, err := converters.TurnToObjectID(expenseID)
	if err != nil {
		return err
	}
	var (
		filter   = bson.D{{Key: "_id", Value: expenseObjectID}, {Key: "is_recurrent", Value: true}}
		updating = bson.D{
			{
				Key:   "$set",
				Value: bson.D{{Key: "is_paid", Value: status}},
			},
		}
	)

	res, err := c.coll.UpdateOne(ctx, filter, updating)
	switch {
	case err != nil:
		return err
	case res.MatchedCount == 0:
		return &repos.NotFoundError{Identifier: expenseID, Entity: "Expense", Message: "it is not recurrent expense"}
	default:
		return nil
	}
}

func (c *ExpensesMongoRepo) FindByNameAndMonthAndIsRecurrent(ctx context.Context, month uint, expenseName string) (*mongoentities.Expense, error) {
	var (
		filter = bson.D{{Key: "name", Value: expenseName}, {Key: "is_recurrent", Value: true}, {Key: "month", Value: month}}
		found  = new(mongoentities.Expense)
	)

	res := c.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		err := res.Err()
		switch {
		case err.Error() == "mongo: no documents in result":
			return nil, &repos.NotFoundError{Identifier: expenseName, Entity: "Recurrent Expense", Message: err.Error()}
		default:
			return nil, err
		}
	}

	if err := res.Decode(&found); err != nil {
		return nil, err
	}

	return found, nil
}

func (c *ExpensesMongoRepo) GetExpenseStatusByID(ctx context.Context, expenseID interface{}) (*schemas.ExpenseIDWithIsPaidStatus, error) {
	expenseObjectID, err := converters.TurnToObjectID(expenseID)
	if err != nil {
		return nil, err
	}
	var (
		filter       = bson.D{{Key: "_id", Value: expenseObjectID}}
		projection   = bson.D{{Key: "_id", Value: 1}, {Key: "is_paid", Value: 1}}
		expenseFound = schemas.ExpenseIDWithIsPaidStatus{}
	)

	if err := c.coll.FindOne(ctx, filter, &options.FindOneOptions{Projection: projection}).Decode(&expenseFound); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &repos.NotFoundError{Identifier: expenseID, Entity: "Expense", Message: "does not exists"}
		}
		return nil, err
	}
	return &expenseFound, nil
}
