package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/gookit/validate"
	"github.com/manicar2093/expenses_api/pkg/period"
	u "github.com/rjNemo/underscore"
	"gopkg.in/guregu/null.v4"
)

type RecurrentExpense struct {
	ID          uuid.UUID          `json:"id,omitempty" gorm:"primaryKey,->"`
	User        *User              `json:"user,omitempty"`
	UserID      uuid.UUID          `json:"user_id,omitempty" validate:"required"`
	Expenses    []*Expense         `json:"expenses,omitempty"`
	Name        string             `json:"name,omitempty" validate:"required"`
	Amount      float64            `json:"amount,omitempty" validate:"required"`
	Description null.String        `json:"description,omitempty"`
	Periodicity period.Periodicity `json:"periodicity" gorm:"<-:update" validate:"required|hasValidData"`
	CreatedAt   *time.Time         `json:"created_at,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty"  gorm:"autoUpdateTime:false"`
}

func (c RecurrentExpense) HasValidData(data period.Periodicity) bool {
	return u.Contains(period.Periods, data)
}

func (c RecurrentExpense) Messages() map[string]string {
	return validate.MS{
		"hasValidData": "periodicity is not into valid types",
	}
}
