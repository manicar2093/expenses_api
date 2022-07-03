package repos

import (
	"context"

	"github.com/go-rel/rel"
	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	IncomesRepository interface {
		Save(context.Context, *entities.Income) error
	}
	IncomesRepositoryImpl struct {
		conn rel.Repository
	}
)

func NewIncomesRepositoryImpl(conn rel.Repository) *IncomesRepositoryImpl {
	return &IncomesRepositoryImpl{
		conn: conn,
	}
}

func (c *IncomesRepositoryImpl) Save(ctx context.Context, income *entities.Income) error {
	return c.conn.Transaction(ctx, func(ctx context.Context) error {
		return c.conn.Insert(ctx, income)
	})
}
