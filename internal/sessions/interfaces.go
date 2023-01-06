package sessions

import (
	"context"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type (
	SessionFindable interface {
		FindByID(ctx context.Context, id uuid.UUID) (*entities.Session, error)
	}
)
