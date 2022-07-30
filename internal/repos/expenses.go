package repos

import (
	"context"
	"time"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ExpensesRepository interface {
		Save(ctx context.Context, expense *entities.Expense) error
		GetExpensesByMonth(ctx context.Context, month time.Month) (*[]entities.Expense, error)
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

func (c *ExpensesRepositoryImpl) GetExpensesByMonth(ctx context.Context, month time.Month) (*[]entities.Expense, error) {
	cursor, err := c.coll.Find(ctx, bson.D{
		{Key: "month", Value: month},
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
