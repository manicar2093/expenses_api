package repos

import (
	"context"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/manicar2093/expenses_api/internal/entities"
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

func (c *RecurrentExpenseGormRepo) Save(ctx context.Context, recurrentExpense *entities.RecurrentExpense) error {
	if res := c.orm.WithContext(ctx).Create(recurrentExpense); res.Error != nil {
		switch err := res.Error.(type) { //nolint:gocritic
		case *pgconn.PgError:
			if strings.Contains(err.Detail, "already exists.") {
				return &AlreadyExistsError{
					Identifier: recurrentExpense.Name,
					Entity:     "Recurrent Expense",
				}
			}
		}
	}
	return nil
}

func (c *RecurrentExpenseGormRepo) FindByName(ctx context.Context, name string) (*entities.RecurrentExpense, error) {
	var found entities.RecurrentExpense
	if res := c.orm.WithContext(ctx).Where("name = ?", name).First(&found); res.Error != nil {
		return nil, res.Error
	}
	return &found, nil
}

func (c *RecurrentExpenseGormRepo) FindAll(ctx context.Context) ([]*entities.RecurrentExpense, error) {
	var found []*entities.RecurrentExpense
	if res := c.orm.Find(&found); res.Error != nil {
		return []*entities.RecurrentExpense{}, res.Error
	}
	return found, nil
}
