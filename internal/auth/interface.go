package auth

import (
	"context"

	"github.com/google/uuid"
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
		Login(ctx context.Context, token string) (*LoginOutput, error)
	}

	TokenRefreshable interface {
		RefreshToken(sessionID uuid.UUID) (string, error)
	}

	OpenIDTokenValidable[T any] interface {
		ValidateOpenIDToken(ctx context.Context, token string) (*T, error)
	}
)
