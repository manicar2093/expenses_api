package expenses

import "github.com/google/uuid"

type (
	ToggleExpenseIsPaidInput struct {
		ID string `json:"id,omitempty" validate:"required|isUUID"`
	}
	ToggleExpenseIsPaidOutput struct {
		ID                  uuid.UUID `json:"id"`
		CurrentIsPaidStatus bool      `json:"current_is_paid_status"`
	}
)
