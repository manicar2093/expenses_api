package expenses

type UpdateExpenseInput struct {
	ID          string  `json:"id,omitempty" validate:"required|isUUID"`
	Name        string  `json:"name,omitempty"`
	Amount      float64 `json:"amount" validate:"required|min:1"`
	Description string  `json:"description,omitempty"`
}
