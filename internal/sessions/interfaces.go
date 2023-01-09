package sessions

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	// SessionValidable checks if stored session is valid with given validation data.
	//
	// This is the package functionality
	SessionValidable interface {
		ValidateSession(ctx context.Context, validationInput *SessionValidationInput) (*entities.Session, error)
	}

	SessionFindable interface {
		FindByID(ctx context.Context, id uuid.UUID) (*entities.Session, error)
	}
)
