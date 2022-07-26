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
	ExpensesRepository interface {
		Save(ctx context.Context, expense *entities.Expense) error
		GetCurrentMonthExpenses(ctx context.Context) (*[]entities.Expense, error)
	}
	ExpensesRepositoryImpl struct {
		coll       *mongo.Collection
		timeGetter dates.TimeGetable
	}
)

func NewExpensesRepositoryImpl(
	coll *mongo.Database,
	timeGetter dates.TimeGetable,
) *ExpensesRepositoryImpl {
	return &ExpensesRepositoryImpl{
		coll:       coll.Collection(entities.ExpenseCollectionName),
		timeGetter: timeGetter,
	}
}

func (c *ExpensesRepositoryImpl) Save(ctx context.Context, expense *entities.Expense) error {
	expense.ID = primitive.NewObjectID()
	createdAt := dates.GetNormalizedDate()
	expense.CreatedAt = &createdAt

	expense.Day = uint(createdAt.Day())
	expense.Month = uint(createdAt.Month())
	expense.Year = uint(createdAt.Year())

	if _, err := c.coll.InsertOne(ctx, expense); err != nil {
		return err
	}

	return nil
}

func (c *ExpensesRepositoryImpl) GetCurrentMonthExpenses(ctx context.Context) (*[]entities.Expense, error) {
	cursor, err := c.coll.Find(ctx, bson.D{
		{Key: "month", Value: c.timeGetter.GetCurrentTime().Month()},
	})
	if err != nil {
		return nil, err
	}

	var response []entities.Expense

	for cursor.Next(ctx) {
		var entityTemp entities.Expense
		if err := cursor.Decode(&entityTemp); err != nil {
			return nil, err
		}
		response = append(response, entityTemp)
	}

	return &response, nil
}
