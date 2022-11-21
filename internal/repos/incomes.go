package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
	"gorm.io/gorm"
)

type (
	IncomesGormRepo struct {
		orm *gorm.DB
	}
)

func NewIncomesGormRepo(conn *gorm.DB) *IncomesGormRepo {
	return &IncomesGormRepo{
		orm: conn,
	}
}

func (c *IncomesGormRepo) Save(ctx context.Context, income *entities.Income) error {
	if res := c.orm.WithContext(ctx).Create(&income); res.Error != nil {
		return res.Error
	}
	return nil
}
