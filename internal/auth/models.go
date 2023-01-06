package auth

import (
	"time"

	"github.com/google/uuid"
)

type (
	LoginOutput struct {
		AccessToken          string    `json:"access_token,omitempty"`
		AccessTokenExpiresAt time.Time `json:"access_token_expires_at,omitempty"`
		RefreshToken         uuid.UUID `json:"refresh_token,omitempty"`
		User                 *UserData `json:"user,omitempty"`
	}

	LoginInput struct {
		Token               string `json:"token"`
		UserAgent, ClientIP string
	}

	UserData struct {
		ID     uuid.UUID `json:"-"`
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
)
