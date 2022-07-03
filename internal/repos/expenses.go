package repos

import (
	"context"

	"github.com/go-rel/rel"
	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	ExpensesRepository interface {
		Save(ctx context.Context, expense *entities.Expense) error
	}
	ExpensesRepositoryImpl struct {
		conn rel.Repository
	}
)

func NewExpensesRepositoryImpl(conn rel.Repository) *ExpensesRepositoryImpl {
	return &ExpensesRepositoryImpl{conn: conn}
}

func (c *ExpensesRepositoryImpl) Save(ctx context.Context, expense *entities.Expense) error {
	return c.conn.Transaction(ctx, func(ctx context.Context) error {
		return c.conn.Insert(ctx, expense)
	})
}
