package repos

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"gorm.io/gorm"
)

type SessionsGormRepo struct {
	orm *gorm.DB
}

func NewSessionsGormRepo(conn *gorm.DB) *SessionsGormRepo {
	return &SessionsGormRepo{orm: conn}
}

func (c *SessionsGormRepo) Save(ctx context.Context, session *entities.Session) error {
	if res := c.orm.WithContext(ctx).Create(session); res.Error != nil {
		return res.Error
	}
	return nil
}

func (c *SessionsGormRepo) FindByID(ctx context.Context, id uuid.UUID) (*entities.Session, error) {
	var found entities.Session
	if res := c.orm.WithContext(ctx).Where("id = ?", id).Preload("User").First(&found); res.Error != nil {
		return nil, res.Error
	}

	return &found, nil
}
