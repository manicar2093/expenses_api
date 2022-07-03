package repos

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/entities"
)

type (
	IncomesRepository interface {
		Save(context.Context, *entities.Income) error
	}
)
