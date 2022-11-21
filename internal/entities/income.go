package entities

import (
	"time"

	"github.com/google/uuid"
)

const IncomeCollectionName = "incomes"

type Income struct {
	ID          uuid.UUID  `json:"id,omitempty" gorm:"primaryKey,->"`
	Name        string     `json:"name,omitempty"`
	Amount      float64    `json:"amount,omitempty"`
	Description string     `json:"description,omitempty"`
	Day         uint       `json:"day,omitempty"`
	Month       uint       `json:"month,omitempty"`
	Year        uint       `json:"year,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
