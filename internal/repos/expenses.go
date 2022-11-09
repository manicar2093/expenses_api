package repos

import (
	"context"
	"errors"
	"time"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/schemas"
	"github.com/manicar2093/expenses_api/pkg/converters"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InsertManyResult struct {
		InsertedIDs []primitive.ObjectID
	}
	ExpensesRepository interface {
		Save(ctx context.Context, expense *entities.Expense) error
		SaveAsRecurrent(ctx context.Context, expense *entities.Expense) error
		GetExpensesByMonth(ctx context.Context, month time.Month) ([]*entities.Expense, error)
		UpdateIsPaidByExpenseID(ctx context.Context, expenseID interface{}, status bool) error
		FindByNameAndMonthAndIsRecurrent(ctx context.Context, month uint, expenseName string) (*entities.Expense, error)
		GetExpenseStatusByID(ctx context.Context, expenseID interface{}) (*schemas.ExpenseIDWithIsPaidStatus, error)
		SaveMany(ctx context.Context, expenses []*entities.Expense) (*InsertManyResult, error)
	}
	ExpensesRepositoryImpl struct {
		coll *mongo.Collection
	}
)

func NewExpensesRepositoryImpl(coll *mongo.Database) *ExpensesRepositoryImpl {
	return &ExpensesRepositoryImpl{coll: coll.Collection(entities.ExpenseCollectionName)}
}

func (c *ExpensesRepositoryImpl) Save(ctx context.Context, expense *entities.Expense) error {
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

func (c *ExpensesRepositoryImpl) GetExpensesByMonth(ctx context.Context, month time.Month) ([]*entities.Expense, error) {
	cursor, err := c.coll.Find(ctx, bson.D{
		{Key: "month", Value: month},
	})
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Expense, 0)

	for cursor.Next(ctx) {
		var entityTemp entities.Expense
		if err := cursor.Decode(&entityTemp); err != nil {
			return nil, err
		}
		response = append(response, &entityTemp)
	}

	return response, nil
}

func (c *ExpensesRepositoryImpl) UpdateIsPaidByExpenseID(ctx context.Context, expenseID interface{}, status bool) error {
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
		return &NotFoundError{Identifier: expenseID, Entity: "Expense", Message: "it is not recurrent expense"}
	default:
		return nil
	}
}

func (c *ExpensesRepositoryImpl) FindByNameAndMonthAndIsRecurrent(ctx context.Context, month uint, expenseName string) (*entities.Expense, error) {
	var (
		filter = bson.D{{Key: "name", Value: expenseName}, {Key: "is_recurrent", Value: true}, {Key: "month", Value: month}}
		found  = new(entities.Expense)
	)

	res := c.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		err := res.Err()
		switch {
		case err.Error() == "mongo: no documents in result":
			return nil, &NotFoundError{Identifier: expenseName, Entity: "Recurrent Expense", Message: err.Error()}
		default:
			return nil, err
		}
	}

	if err := res.Decode(&found); err != nil {
		return nil, err
	}

	return found, nil
}

func (c *ExpensesRepositoryImpl) GetExpenseStatusByID(ctx context.Context, expenseID interface{}) (*schemas.ExpenseIDWithIsPaidStatus, error) {
	expenseObjectID, err := converters.TurnToObjectID(expenseID)
	if err != nil {
		return nil, err
	}
	var (
		filter     = bson.D{{Key: "_id", Value: expenseObjectID}}
		projection = bson.D{
			{Key: "_id", Value: 1},
			{Key: "is_paid", Value: 1},
			{Key: "recurrent_expense_id", Value: 1},
		}
		expenseFound = schemas.ExpenseIDWithIsPaidStatus{}
	)

	if err := c.coll.FindOne(ctx, filter, &options.FindOneOptions{Projection: projection}).Decode(&expenseFound); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &NotFoundError{Identifier: expenseID, Entity: "Expense", Message: "does not exists"}
		}
		return nil, err
	}
	return &expenseFound, nil
}

func (c *ExpensesRepositoryImpl) SaveAsRecurrent(ctx context.Context, expense *entities.Expense) error {
	expense.IsRecurrent = true
	return c.Save(ctx, expense)
}

func (c *ExpensesRepositoryImpl) SaveMany(
	ctx context.Context,
	expenses []*entities.Expense,
) (*InsertManyResult, error) {
	insertable := []interface{}{}
	for _, expense := range expenses {
		expense.ID = primitive.NewObjectID()
		insertable = append(insertable, expense)
	}

	got, err := c.coll.InsertMany(ctx, insertable)
	if err != nil {
		return nil, err
	}

	result := &InsertManyResult{}
	for _, id := range got.InsertedIDs {
		result.InsertedIDs = append(result.InsertedIDs, id.(primitive.ObjectID))
	}

	return result, nil
}
