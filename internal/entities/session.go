package entities

import (
	"time"

	"github.com/google/uuid"
)

const SessionsTable = "sessions"

type Session struct {
	ID        uuid.UUID  `json:"id,omitempty"`
	User      *User      `json:"user,omitempty"`
	UserID    uuid.UUID  `json:"user_id,omitempty"`
	UserAgent string     `json:"user_agent,omitempty"`
	ClientIP  string     `json:"client_ip,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
