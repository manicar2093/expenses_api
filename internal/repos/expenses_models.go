package repos

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type UpdateExpenseInput struct {
	ID          uuid.UUID
	Name        null.String
	Amount      float64
	Description null.String
}

func (c *UpdateExpenseInput) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":        c.Name,
		"amount":      c.Amount,
		"description": c.Description,
	}
}
