package repos_test

import (
	"context"
	"math"
	"net/http"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"
)

var _ = Describe("Expenses", func() {
	var (
		ctx                            context.Context
		repo                           repos.ExpensesRepository
		expectedUserID                 uuid.UUID
		expectedUser                   entities.User
		expectedRecurrentExpenseID     uuid.UUID
		expectedRecurrentExpenseNullID uuid.NullUUID
		expectedRecurrentExpense       entities.RecurrentExpense
	)

	BeforeEach(func() {
		ctx = context.TODO()
		repo = repos.NewExpensesGormRepo(conn)
		expectedUserID = uuid.New()
		expectedUser = entities.User{
			ID:       expectedUserID,
			Name:     null.NewString(faker.Name(), true),
			Lastname: null.NewString(faker.LastName(), true),
			Email:    faker.Email(),
		}
		expectedRecurrentExpenseID = uuid.New()
		expectedRecurrentExpense = entities.RecurrentExpense{
			ID:          expectedRecurrentExpenseID,
			UserID:      expectedUserID,
			Name:        faker.Name(),
			Amount:      faker.Latitude(),
			Description: null.StringFrom(faker.Paragraph()),
		}
		expectedRecurrentExpenseNullID = uuid.NullUUID{
			UUID:  expectedRecurrentExpenseID,
			Valid: true,
		}
		conn.Create(&expectedUser)
		conn.Create(&expectedRecurrentExpense)
	})

	AfterEach(func() {
		conn.Delete(&expectedRecurrentExpense)
		conn.Delete(&expectedUser)
	})

	Describe("Save", func() {
		It("saves an entities.Expense in database", func() {

			var (
				expectedName        = null.StringFrom(faker.Name())
				expectedAmount      = faker.Latitude()
				expectedDescription = faker.Sentence()
				expectedExpense     = entities.Expense{
					Name:               expectedName,
					UserID:             expectedUserID,
					RecurrentExpenseID: expectedRecurrentExpenseNullID,
					Amount:             expectedAmount,
					Description:        null.StringFrom(expectedDescription),
				}
			)
			defer conn.Delete(&expectedExpense)

			err := repo.Save(ctx, &expectedExpense)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedExpense.ID.String()).ToNot(BeEmpty())
			Expect(expectedExpense.IsPaid).To(BeFalse())
			Expect(expectedExpense.CreatedAt).ToNot(BeZero())
			Expect(expectedExpense.UpdatedAt).To(BeZero())

		})
	})

	Describe("UpdateIsPaidByExpenseID", func() {
		It("change isPaid by given bool", func() {
			var (
				expectedName        = null.StringFrom(faker.Name())
				expectedAmount      = faker.Latitude()
				expectedDescription = faker.Sentence()
				expectedStatus      = true
				expectedExpense     = &entities.Expense{
					Name:               expectedName,
					UserID:             expectedUserID,
					RecurrentExpenseID: expectedRecurrentExpenseNullID,
					Amount:             expectedAmount,
					Description:        null.StringFrom(expectedDescription),
				}
			)
			conn.Create(&expectedExpense)
			defer conn.Delete(&expectedExpense)

			err := repo.UpdateIsPaidByExpenseID(ctx, expectedExpense.ID, expectedStatus)

			var changed entities.Expense
			conn.Find(&changed, "id = ?", expectedExpense.ID)

			Expect(err).ToNot(HaveOccurred())
			Expect(changed.Name).To(Equal(expectedName))
			Expect(changed.IsPaid).To(Equal(expectedStatus))
		})

		When("expense is not found", func() {
			It("returns a NotFound error", func() {
				var (
					expectedID = uuid.New()
				)

				err := repo.UpdateIsPaidByExpenseID(ctx, expectedID, true)

				Expect(err).To(BeAssignableToTypeOf(&apperrors.NotFoundError{}))
				Expect(err.(*apperrors.NotFoundError).StatusCode()).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("GetExpenseStatusByID", func() {
		It("finds a expense by ID retriving just is_paid and its ID", func() {
			var (
				expectedName        = null.StringFrom(faker.Name())
				expectedAmount      = faker.Latitude()
				expectedDescription = faker.Sentence()
				expectedStatus      = true
				expectedExpense     = &entities.Expense{
					Name:        expectedName,
					UserID:      expectedUserID,
					Amount:      expectedAmount,
					IsPaid:      expectedStatus,
					Description: null.StringFrom(expectedDescription),
				}
			)
			conn.Create(&expectedExpense)
			defer conn.Delete(&expectedExpense)

			got, err := repo.GetExpenseStatusByID(ctx, expectedExpense.ID)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.ID).To(Equal(expectedExpense.ID))
			Expect(got.IsPaid).To(Equal(expectedStatus))

		})

		When("expense is not found", func() {
			It("returns a NotFound error", func() {
				var (
					expectedID = uuid.New()
				)

				got, err := repo.GetExpenseStatusByID(ctx, expectedID)

				Expect(got).To(BeNil())
				Expect(err).To(BeAssignableToTypeOf(&apperrors.NotFoundError{}))
				Expect(err.(*apperrors.NotFoundError).StatusCode()).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("Update", func() {
		It("changes expense required saved data", func() {
			var (
				savedExpenseID = uuid.New()
				savedExpense   = entities.Expense{
					ID:          savedExpenseID,
					UserID:      expectedUserID,
					Name:        null.StringFrom(faker.Name()),
					Amount:      faker.Latitude(),
					Description: null.StringFrom(faker.Paragraph()),
					IsPaid:      true,
				}
				expectedNewName             = null.StringFrom(faker.Name())
				expectedNewAmount           = math.Round(faker.Latitude())
				expectedNewDescription      = null.StringFrom(faker.Name())
				expectedExpenseDataToUpdate = repos.UpdateExpenseInput{
					ID:          savedExpenseID,
					Name:        expectedNewName,
					Amount:      expectedNewAmount,
					Description: expectedNewDescription,
				}
			)
			conn.Create(&savedExpense)
			defer conn.Delete(&savedExpense)

			err := repo.Update(ctx, &expectedExpenseDataToUpdate)
			var updated *entities.Expense
			conn.Where("id = ?", savedExpense.ID).Find(&updated)

			Expect(err).ToNot(HaveOccurred())
			Expect(updated.Name).To(Equal(expectedNewName))
			Expect(updated.Amount).To(Equal(expectedNewAmount))
			Expect(updated.Description).To(Equal(expectedNewDescription))
			Expect(updated.IsPaid).To(Equal(savedExpense.IsPaid))
		})

		When("expense does not exists", func() {
			It("return a notFoundError", func() {
				var (
					savedExpenseID              = uuid.New()
					expectedExpenseDataToUpdate = repos.UpdateExpenseInput{
						ID:          savedExpenseID,
						Name:        null.StringFrom(faker.Name()),
						Amount:      faker.Latitude(),
						Description: null.StringFrom(faker.Paragraph()),
					}
				)

				err := repo.Update(ctx, &expectedExpenseDataToUpdate)

				Expect(err).To(BeAssignableToTypeOf(&apperrors.NotFoundError{}))
			})
		})
	})

	Describe("FindByID", func() {
		It("returns an expense found by its ID", func() {
			var (
				expectedExpenseID = uuid.New()
				savedExpense      = entities.Expense{
					ID:                 expectedExpenseID,
					UserID:             expectedUserID,
					RecurrentExpenseID: expectedRecurrentExpenseNullID,
					Name:               null.StringFrom(faker.Name()),
					Amount:             faker.Latitude(),
				}
			)
			conn.Create(&savedExpense)
			defer conn.Delete(&savedExpense)

			got, err := repo.FindByID(ctx, expectedExpenseID)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.ID).To(Equal(expectedExpenseID))
			Expect(got.RecurrentExpenseID).To(Equal(expectedRecurrentExpenseNullID))
			Expect(got.RecurrentExpense.ID).To(Equal(expectedRecurrentExpenseID))
		})

		When("expense does not exists", func() {
			It("return a notFoundError", func() {
				var (
					savedExpenseID = uuid.New()
				)

				got, err := repo.FindByID(ctx, savedExpenseID)

				Expect(err).To(BeAssignableToTypeOf(&apperrors.NotFoundError{}))
				Expect(got).To(BeNil())
			})
		})
	})
})
