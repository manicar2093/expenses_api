package auth

import (
	"time"

	"github.com/google/uuid"
)

type (
	LoginOutput struct {
		SessionID             uuid.UUID `json:"session_id,omitempty"`
		AccessToken           string    `json:"access_token,omitempty"`
		AccessTokenExpiresAt  time.Time `json:"access_token_expires_at,omitempty"`
		RefreshToken          string    `json:"refresh_token,omitempty"`
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
		User                  *UserData `json:"user,omitempty"`
	}

	UserData struct {
		Name   string `json:"name,omitempty"`
		Email  string `json:"email,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}

	AccessToken struct {
		Expiration time.Duration `json:"expiration,omitempty"`
		UserID     uuid.UUID     `json:"user_id,omitempty"`
	}

	RefreshToken struct {
		Expiration time.Duration `json:"expiration,omitempty"`
		SessionID  uuid.UUID     `json:"user_id,omitempty"`
	}

	Tokenizable interface {
		CreateAccessToken(tokenDetails *AccessToken) (string, error)
		CreateRefreshToken(tokenDetails *RefreshToken) (string, error)
		ValidateToken(token string) error
	}

	GoogleLoginable interface {
		Login(googleSignInToken string) (*LoginOutput, error)
	}

	TokenRefreshable interface {
		RefreshToken(refreshToken string) (string, error)
	}
)
