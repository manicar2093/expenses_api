package repos

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
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

func (c *ExpensesGormRepo) GetExpensesByMonth(ctx context.Context, month time.Month, userID uuid.UUID) ([]*entities.Expense, error) {
	var expensesFound []*entities.Expense
	if res := c.orm.WithContext(ctx).Where(&entities.Expense{Month: uint(month), UserID: userID}, month).Preload("RecurrentExpense").Find(&expensesFound); res.Error != nil {
		return []*entities.Expense{}, res.Error
	}
	return expensesFound, nil
}

func (c *ExpensesGormRepo) UpdateIsPaidByExpenseID(ctx context.Context, expenseID uuid.UUID, status bool) error {
	res := c.orm.WithContext(ctx).Model(&entities.Expense{}).Where("id = ?", expenseID).Update("is_paid", status)
	switch {
	case res.Error != nil:
		return res.Error
	case res.RowsAffected == 0:
		err := &apperrors.NotFoundError{Identifier: expenseID, Entity: entities.ExpensesEntityName, Message: "can´t be updated. It does not exist"}
		log.Errorln(err)
		return err
	}
	return nil
}

func (c *ExpensesGormRepo) FindByNameAndMonthAndIsRecurrent(
	ctx context.Context,
	month uint,
	expenseName string,
	userID uuid.UUID,
) (*entities.Expense, error) {
	var found entities.Expense
	if res := c.orm.WithContext(ctx).Where("user_id = ? AND month = ? AND name = ? AND recurrent_expense_id IS NOT null", userID, month, expenseName).First(&found); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, &apperrors.NotFoundError{Identifier: expenseName, Entity: entities.ExpensesEntityName, Message: res.Error.Error()}
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
		return nil, &apperrors.NotFoundError{Identifier: expenseID, Entity: entities.ExpensesEntityName, Message: "Any row found with data"}
	}
	return &found, nil
}

func (c *ExpensesGormRepo) Update(ctx context.Context, expenseUpdateInput *UpdateExpenseInput) error {
	res := c.orm.WithContext(ctx).Model(
		&entities.Expense{ID: expenseUpdateInput.ID},
	).Select("Name", "Amount", "Description").Updates(expenseUpdateInput.ToMap())
	switch {
	case res.Error != nil:
		return res.Error
	case res.RowsAffected == 0:
		return &apperrors.NotFoundError{Identifier: expenseUpdateInput.ID, Entity: entities.ExpensesEntityName, Message: "can't update. It is not in db"}
	}
	return nil
}

func (c *ExpensesGormRepo) FindByID(ctx context.Context, expenseID uuid.UUID) (*entities.Expense, error) {
	var found entities.Expense
	if res := c.orm.WithContext(ctx).Preload("RecurrentExpense").First(&found, "id = ?", expenseID); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, &apperrors.NotFoundError{Identifier: expenseID, Entity: entities.ExpensesEntityName, Message: res.Error.Error()}
		}
		return nil, res.Error
	}
	return &found, nil
}
