package repos

import (
	"context"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"gorm.io/gorm"
)

type (
	RecurrentExpenseGormRepo struct {
		orm *gorm.DB
	}
)

func NewRecurrentExpenseGormRepo(conn *gorm.DB) *RecurrentExpenseGormRepo {
	return &RecurrentExpenseGormRepo{
		orm: conn,
	}
}

func (c *RecurrentExpenseGormRepo) Create(ctx context.Context, recurrentExpense *entities.RecurrentExpense) error {
	if res := c.orm.WithContext(ctx).Create(recurrentExpense); res.Error != nil {
		switch err := res.Error.(type) { //nolint:gocritic
		case *pgconn.PgError:
			if strings.Contains(err.Detail, "already exists.") {
				return &apperrors.AlreadyExistsError{
					Identifier: recurrentExpense.Name,
					Entity:     "Recurrent Expense",
				}
			}
		default:
			return res.Error
		}
	}
	return nil
}
