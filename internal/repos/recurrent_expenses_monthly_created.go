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
	RecurrentExpensesMonthlyCreatedRepo interface {
		Save(
			ctx context.Context,
			recurrentExpense *entities.RecurrentExpensesMonthlyCreated,
		) error
		FindByCurrentMonthAndYear(ctx context.Context, month uint, year uint) (*entities.RecurrentExpensesMonthlyCreated, error)
		Update(ctx context.Context, recurrentExpense *entities.RecurrentExpensesMonthlyCreated) error
	}
	RecurrentExpensesMonthlyCreatedRepoImpl struct {
		recurrentExpensesColl, coll *mongo.Collection
	}
)

func NewRecurrentExpensesCreatedMonthlyRepoImpl(
	conn *mongo.Database,
) *RecurrentExpensesMonthlyCreatedRepoImpl {
	return &RecurrentExpensesMonthlyCreatedRepoImpl{
		coll:                  conn.Collection(entities.RecurrentExpensesCreatedMonthlyCollectionName),
		recurrentExpensesColl: conn.Collection(entities.RecurrentExpensesCollectonName),
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

func (c *RecurrentExpensesMonthlyCreatedRepoImpl) FindByCurrentMonthAndYear(ctx context.Context, month uint, year uint) (*entities.RecurrentExpensesMonthlyCreated, error) {
	var (
		found entities.RecurrentExpensesMonthlyCreated
	)

	cursor, err := c.coll.Aggregate(ctx, bson.A{
		bson.D{
			{
				Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: entities.ExpensesCountCollName},
					{Key: "localField", Value: "_id"},
					{Key: "foreignField", Value: "recurrent_expenses_monthly_created_id"},
					{Key: "as", Value: "expenses_count"},
				},
			},
		},
		bson.D{
			{
				Key: "$match",
				Value: bson.D{
					{Key: "month", Value: month},
					{Key: "year", Value: year},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if !cursor.Next(ctx) {
		return nil, &NotFoundError{Identifier: map[string]interface{}{
			"month": &month,
			"year":  &year,
		}, Entity: "RecurrentExpensesMonthlyCreated", Message: "no results for aggregation operation"}
	}
	if err := cursor.Decode(&found); err != nil {
		return nil, err
	}
	for index, data := range found.ExpensesCount {
		var temp entities.RecurrentExpense
		if err := c.recurrentExpensesColl.FindOne(ctx, bson.D{{Key: "_id", Value: data.RecurrentExpenseID}}).Decode(&temp); err != nil {
			return nil, err
		}
		found.ExpensesCount[index].RecurrentExpense = &temp
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
