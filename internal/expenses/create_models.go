package expenses

type CreateExpenseInput struct {
	Name         string  `json:"name,omitempty" validate:"required"`
	Amount       float64 `json:"amount,omitempty" validate:"required"`
	Description  string  `json:"description,omitempty" validate:"-"`
	ForNextMonth bool    `json:"for_next_month,omitempty" validate:"-"`
}
