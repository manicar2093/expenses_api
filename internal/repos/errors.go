package repos

import (
	"errors"

	"gorm.io/gorm"
)

func isNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
