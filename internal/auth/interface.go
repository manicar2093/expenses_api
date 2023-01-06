package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	UserAuthenticable interface {
		FindUserByEmail(ctx context.Context, email string) (*UserData, error)
		CreateUser(ctx context.Context, user *UserData) error
	}

	Tokenizable interface {
		CreateAccessToken(tokenDetails *AccessToken) (*TokenInfo, error)
		CreateRefreshToken(tokenDetails *RefreshToken) (*TokenInfo, error)
	}
	TokenValidable interface {
		ValidateToken(ctx context.Context, token string, output interface{}) error
	}

	LoginableByToken interface {
		Login(ctx context.Context, loginInput *LoginInput) (*LoginOutput, error)
	}

	TokenRefreshable interface {
		RefreshToken(sessionID uuid.UUID) (*LoginOutput, error)
	}

	SessionCreateable interface {
		Create(ctx context.Context, session *entities.Session) error
	}

	OpenIDTokenValidable[T any] interface {
		ValidateOpenIDToken(ctx context.Context, token string) (*T, error)
	}
)
