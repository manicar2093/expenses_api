package sessions

import "github.com/google/uuid"

type SessionValidationInput struct {
	SessionID     uuid.UUID
	FromUserAgent string
	FromClientIP  string
}
