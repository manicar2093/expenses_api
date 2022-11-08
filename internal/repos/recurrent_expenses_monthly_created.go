package repos

import (
	"context"
	"errors"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	RecurrentExpensesMonthlyCreatedRepo interface {
		Save(
			ctx context.Context,
			recurrentExpense *entities.RecurrentExpensesMonthlyCreated,
		) error
		FindByMonthAndYear(ctx context.Context, month uint, year uint) (*entities.RecurrentExpensesMonthlyCreated, error)
		Update(ctx context.Context, recurrentExpense *entities.RecurrentExpensesMonthlyCreated) error
	}
	RecurrentExpensesMonthlyCreatedRepoImpl struct {
		coll *mongo.Collection
	}
)

func NewRecurrentExpensesCreatedMonthlyRepoImpl(
	conn *mongo.Database,
) *RecurrentExpensesMonthlyCreatedRepoImpl {
	return &RecurrentExpensesMonthlyCreatedRepoImpl{
		coll: conn.Collection(entities.RecurrentExpensesCreatedMonthlyCollectionName),
	}
}

func (c *RecurrentExpensesMonthlyCreatedRepoImpl) Save(
	ctx context.Context,
	recurrentExpense *entities.RecurrentExpensesMonthlyCreated,
) error {
	nomalizedDate := dates.GetNormalizedDate()
	recurrentExpense.ID = primitive.NewObjectID()
	recurrentExpense.CreatedAt = &nomalizedDate
	_, err := c.coll.InsertOne(ctx, recurrentExpense)
	if err != nil {
		return err
	}
	return nil
}

func (c *RecurrentExpensesMonthlyCreatedRepoImpl) FindByMonthAndYear(ctx context.Context, month uint, year uint) (*entities.RecurrentExpensesMonthlyCreated, error) {
	var (
		found   entities.RecurrentExpensesMonthlyCreated
		filters = primitive.D{
			{Key: "month", Value: month},
			{Key: "year", Value: year},
		}
	)
	if err := c.coll.FindOne(ctx, filters).Decode(&found); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &NotFoundError{Identifier: map[string]interface{}{
				"month": &month,
				"year":  &year,
			}, Entity: "RecurrentExpensesMonthlyCreated", Message: err.Error()}
		}
		return nil, err
	}

	return &found, nil
}

func (c *RecurrentExpensesMonthlyCreatedRepoImpl) Update(ctx context.Context, recurrentExpense *entities.RecurrentExpensesMonthlyCreated) error {
	_, err := c.coll.ReplaceOne(ctx, bson.D{{Key: "_id", Value: recurrentExpense.ID}}, recurrentExpense)
	if err != nil {
		return err
	}
	return nil
}
