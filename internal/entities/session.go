package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id,omitempty" gorm:"primaryKey,->"`
	User         *User     `json:"user,omitempty"`
	UserID       uuid.UUID `json:"user_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	ClientIP     string    `json:"client_ip,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}
