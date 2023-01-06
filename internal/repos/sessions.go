package repos

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"gorm.io/gorm"
)

type SessionsGormRepo struct {
	orm *gorm.DB
}

func NewSessionGormRepo(conn *gorm.DB) *SessionsGormRepo {
	return &SessionsGormRepo{orm: conn}
}

func (c *SessionsGormRepo) FindByID(ctx context.Context, id uuid.UUID) (*entities.Session, error) {
	var found entities.Session
	if res := c.orm.WithContext(ctx).Table(entities.SessionsTable).Where("id = ?", id).First(&found); res.Error != nil {
		if isNotFoundError(res.Error) {
			return nil, &apperrors.NotFoundError{Identifier: id, Entity: "Session", Message: res.Error.Error()}
		}
		return nil, res.Error
	}
	return &found, nil
}

func (c *SessionsGormRepo) Create(ctx context.Context, session *entities.Session) error {
	if res := c.orm.WithContext(ctx).Create(session); res.Error != nil {
		return res.Error
	}
	return nil
}
