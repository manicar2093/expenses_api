package entities

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

const UserTable = "users"

type User struct {
	ID        uuid.UUID   `json:"id,omitempty" gorm:"primaryKey,->"`
	Name      null.String `json:"name,omitempty"`
	Lastname  null.String `json:"lastname,omitempty"`
	Email     string      `json:"email,omitempty"`
	Avatar    null.String `json:"avatar,omitempty"`
	CreatedAt *time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time  `json:"updated_at,omitempty"`
}
