package sessions

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
)

type DefaultValidator struct {
	sessionFindable SessionFindable
}

func NewDefaultValidator(sessionFindable SessionFindable) *DefaultValidator {
	return &DefaultValidator{sessionFindable: sessionFindable}
}

func (c *DefaultValidator) ValidateSession(ctx context.Context, validationInput *SessionValidationInput) (*entities.Session, error) {
	session, err := c.sessionFindable.FindByID(ctx, validationInput.SessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}
