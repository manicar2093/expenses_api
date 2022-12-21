package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	LoginOutput struct {
		AccessToken           string    `json:"access_token,omitempty"`
		AccessTokenExpiresAt  time.Time `json:"access_token_expires_at,omitempty"`
		RefreshToken          string    `json:"refresh_token,omitempty"`
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
		User                  *UserData `json:"user,omitempty"`
	}

	UserData struct {
		ID     uuid.UUID `json:"-,omitempty"`
		Name   string    `json:"name,omitempty"`
		Email  string    `json:"email,omitempty"`
		Avatar string    `json:"avatar,omitempty"`
	}

	AccessToken struct {
		Expiration time.Duration `json:"expiration,omitempty"`
		UserID     uuid.UUID     `json:"user_id,omitempty"`
	}

	RefreshToken struct {
		Expiration time.Duration `json:"expiration,omitempty"`
		SessionID  uuid.UUID     `json:"user_id,omitempty"`
	}

	TokenInfo struct {
		Token     string
		ExpiresAt time.Time
	}

	UserAuthenticable interface {
		FindUserByEmail(ctx context.Context, email string) (*UserData, error)
		CreateUser(ctx context.Context, user *UserData) error
	}

	Tokenizable interface {
		CreateAccessToken(tokenDetails *AccessToken) (*TokenInfo, error)
		CreateRefreshToken(tokenDetails *RefreshToken) (*TokenInfo, error)
		ValidateToken(token string) error
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
