package repos

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"gorm.io/gorm"
)

type ExpensesGormRepo struct {
	orm *gorm.DB
}

func NewExpensesGormRepo(orm *gorm.DB) *ExpensesGormRepo {
	return &ExpensesGormRepo{orm: orm}
}

func (c *ExpensesGormRepo) Save(ctx context.Context, expense *entities.Expense) error {
	if res := c.orm.WithContext(ctx).Create(expense); res.Error != nil {
		return res.Error
	}
	return nil
}

func (c *ExpensesGormRepo) GetExpensesByMonth(ctx context.Context, month time.Month) ([]*entities.Expense, error) {
	var expensesFound []*entities.Expense
	if res := c.orm.WithContext(ctx).Where(&entities.Expense{Month: uint(month)}, month).Find(&expensesFound); res.Error != nil {
		return []*entities.Expense{}, res.Error
	}
	return expensesFound, nil
}

func (c *ExpensesGormRepo) UpdateIsPaidByExpenseID(ctx context.Context, expenseID uuid.UUID, status bool) error {
	if res := c.orm.WithContext(ctx).Model(&entities.Expense{}).Where("id = ?", expenseID).Update("is_paid", status); res.Error != nil {
		return res.Error
	}
	return nil
}

func (c *ExpensesGormRepo) FindByNameAndMonthAndIsRecurrent(ctx context.Context, month uint, expenseName string) (*entities.Expense, error) {
	var found entities.Expense
	if res := c.orm.WithContext(ctx).Where("month = ? AND name = ? AND recurrent_expense_id IS NOT null", month, expenseName).First(&found); res.Error != nil {
		switch {
		case errors.Is(res.Error, gorm.ErrRecordNotFound):
			return nil, &NotFoundError{Identifier: expenseName, Entity: "Expense", Message: res.Error.Error()}
		}
		return nil, res.Error
	}

	return &found, nil
}

func (c *ExpensesGormRepo) GetExpenseStatusByID(ctx context.Context, expenseID uuid.UUID) (*entities.ExpenseIDWithIsPaidStatus, error) {
	var found entities.ExpenseIDWithIsPaidStatus
	res := c.orm.WithContext(ctx).Table(
		entities.ExpensesTableName,
	).Select("id", "is_paid").Where(
		"id = ?", expenseID,
	).Scan(&found)
	switch {
	case res.Error != nil:
		return nil, res.Error
	case res.RowsAffected == 0:
		return nil, &NotFoundError{Identifier: expenseID, Entity: "Expense", Message: "Any row found with data"}
	}
	return &found, nil
}
