package repos

import (
	"context"
	"errors"

	"github.com/manicar2093/expenses_api/internal/entities"
	"gorm.io/gorm"
)

type UserGormRepo struct {
	orm *gorm.DB
}

func NewUserGormRepo(conn *gorm.DB) *UserGormRepo {
	return &UserGormRepo{orm: conn}
}

func (c *UserGormRepo) Save(ctx context.Context, user *entities.User) error {
	return c.orm.WithContext(ctx).Create(&user).Error
}

func (c *UserGormRepo) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var found entities.User
	if res := c.orm.WithContext(ctx).Where("email = ?", email).First(&found); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, &NotFoundError{Identifier: email, Entity: "User", Message: res.Error.Error()}
		}
		return nil, res.Error
	}

	return &found, nil
}
