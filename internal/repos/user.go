package repos

import (
	"context"
	"errors"

	"github.com/manicar2093/expenses_api/internal/auth"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/nullsql"
	"gorm.io/gorm"
)

type UserGormRepo struct {
	orm *gorm.DB
}

func NewUserGormRepo(conn *gorm.DB) *UserGormRepo {
	return &UserGormRepo{orm: conn}
}

func (c *UserGormRepo) CreateUser(ctx context.Context, user *auth.UserData) error {
	var userToCreate = entities.User{
		Name:   nullsql.ValidateStringSQLValid(user.Name),
		Email:  user.Email,
		Avatar: nullsql.ValidateStringSQLValid(user.Avatar),
	}
	if res := c.orm.WithContext(ctx).Create(&userToCreate); res.Error != nil {
		return res.Error
	}
	user.ID = userToCreate.ID
	return nil
}

func (c *UserGormRepo) FindUserByEmail(ctx context.Context, email string) (*auth.UserData, error) {
	var found auth.UserData
	if res := c.orm.WithContext(ctx).Table(entities.UserTable).Where("email = ?", email).First(&found); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, &apperrors.NotFoundError{Identifier: email, Entity: "User", Message: res.Error.Error()}
		}
		return nil, res.Error
	}

	return &found, nil
}
