package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/manicar2093/goption"
)

const IncomeCollectionName = "incomes"

type Income struct {
	ID          uuid.UUID                   `json:"id,omitempty" gorm:"primaryKey,->"`
	UserID      goption.Optional[uuid.UUID] `json:"user_id,omitempty" validate:"required"`
	Name        string                      `json:"name,omitempty" validate:"required"`
	Amount      float64                     `json:"amount,omitempty" validate:"required"`
	Description string                      `json:"description,omitempty"`
	CreatedAt   time.Time                   `json:"created_at,omitempty"`
	UpdatedAt   goption.Optional[time.Time] `json:"updated_at,omitempty" gorm:"autoUpdateTime:false"`
}
