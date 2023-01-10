package incomes

import "github.com/manicar2093/expenses_api/internal/entities"

type CreateIncomeInput struct {
	entities.Income `json:",inline" validate:"required"`
}
