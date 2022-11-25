package expenses

type SetExpenseToPaidInput struct {
	ID string `json:"id,omitempty" validate:"required|isUUID"`
}
