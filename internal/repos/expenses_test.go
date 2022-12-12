package repos_test

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"
)

var _ = Describe("Expenses", func() {
	var (
		ctx                            context.Context
		repo                           *repos.ExpensesGormRepo
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
					Day:                1,
					Month:              1,
					Year:               2022,
					Description:        null.StringFrom(expectedDescription),
				}
			)
			defer conn.Delete(&expectedExpense)

			err := repo.Save(ctx, &expectedExpense)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedExpense.ID.String()).ToNot(BeEmpty())
			Expect(expectedExpense.Day).ToNot(BeZero())
			Expect(expectedExpense.Month).ToNot(BeZero())
			Expect(expectedExpense.Year).ToNot(BeZero())
			Expect(expectedExpense.IsPaid).To(BeFalse())
			Expect(expectedExpense.CreatedAt).ToNot(BeZero())
			Expect(expectedExpense.UpdatedAt).ToNot(BeZero())

		})
	})

	Describe("GetExpensesByMonth", func() {
		It("returns all expenses by current month", func() {
			var (
				expectedName        = null.StringFrom(faker.Name())
				expectedAmount      = faker.Latitude()
				expectedDescription = faker.Sentence()
				expectedMonth       = time.January
				expectedExpenses    = []*entities.Expense{
					{
						Name:               expectedName,
						UserID:             expectedUserID,
						RecurrentExpenseID: expectedRecurrentExpenseNullID,
						Amount:             expectedAmount,
						Day:                1,
						Month:              uint(expectedMonth),
						Year:               2022,
						Description:        null.StringFrom(expectedDescription),
					},
					{
						Name:               expectedName,
						UserID:             expectedUserID,
						RecurrentExpenseID: expectedRecurrentExpenseNullID,
						Amount:             expectedAmount,
						Day:                1,
						Month:              uint(expectedMonth),
						Year:               2022,
						Description:        null.StringFrom(expectedDescription),
					},
					{
						Name:               expectedName,
						UserID:             expectedUserID,
						RecurrentExpenseID: expectedRecurrentExpenseNullID,
						Amount:             expectedAmount,
						Day:                1,
						Month:              uint(expectedMonth),
						Year:               2022,
						Description:        null.StringFrom(expectedDescription),
					},
				}
			)

			conn.Create(&expectedExpenses)
			defer conn.Delete(&expectedExpenses)

			got, err := repo.GetExpensesByMonth(ctx, expectedMonth)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(HaveLen(3))
			Expect(got[0]).To(BeAssignableToTypeOf(&entities.Expense{}))
		})

		When("There is no data saved", func() {
			It("returns an empty slice", func() {
				got, err := repo.GetExpensesByMonth(ctx, time.July)

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(HaveLen(0))
			})
		})
	})

	Describe("UpdateIsPaidByExpenseID", func() {
		It("change isPaid by given bool", func() {
			var (
				expectedName        = null.StringFrom(faker.Name())
				expectedAmount      = faker.Latitude()
				expectedDescription = faker.Sentence()
				expectedMonth       = time.January
				expectedStatus      = true
				expectedExpense     = &entities.Expense{
					Name:               expectedName,
					UserID:             expectedUserID,
					RecurrentExpenseID: expectedRecurrentExpenseNullID,
					Amount:             expectedAmount,
					Day:                1,
					Month:              uint(expectedMonth),
					Year:               2022,
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

				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(err.(*repos.NotFoundError).StatusCode()).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("FindByNameAndMonthAndIsRecurrent", func() {
		It("finds a expense that is recurrent by its name", func() {
			var (
				expectedRecurrenteExpenseName = faker.Name()
				expectedAmount                = faker.Latitude()
				expectedDescription           = faker.Sentence()
				expectedMonth                 = time.January
				expectedExpenseName           = null.StringFrom(expectedRecurrenteExpenseName)
				expectedExpense               = &entities.Expense{
					Name:               expectedExpenseName,
					UserID:             expectedUserID,
					RecurrentExpenseID: expectedRecurrentExpenseNullID,
					Amount:             expectedAmount,
					Day:                1,
					Month:              uint(expectedMonth),
					Year:               2022,
					Description:        null.StringFrom(expectedDescription),
				}
			)
			conn.Create(&expectedExpense)
			defer conn.Delete(&expectedExpense)

			got, err := repo.FindByNameAndMonthAndIsRecurrent(ctx, uint(expectedMonth), expectedRecurrenteExpenseName)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.Name).To(Equal(expectedExpenseName))
			Expect(got.RecurrentExpenseID).ToNot(BeNil())
		})

		When("expense does not exist", func() {
			It("return an repos.NotFoundError", func() {
				var (
					expectedRecurrenteExpenseName = faker.Name()
					expectedMonth                 = time.January
				)

				got, err := repo.FindByNameAndMonthAndIsRecurrent(ctx, uint(expectedMonth), expectedRecurrenteExpenseName)

				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(got).To(BeNil())
			})
		})
	})

	Describe("GetExpenseStatusByID", func() {
		It("finds a expense by ID retriving just is_paid and its ID", func() {
			var (
				expectedName        = null.StringFrom(faker.Name())
				expectedAmount      = faker.Latitude()
				expectedDescription = faker.Sentence()
				expectedMonth       = time.January
				expectedStatus      = true
				expectedExpense     = &entities.Expense{
					Name:        expectedName,
					UserID:      expectedUserID,
					Amount:      expectedAmount,
					Day:         1,
					Month:       uint(expectedMonth),
					Year:        2022,
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
				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(err.(*repos.NotFoundError).StatusCode()).To(Equal(http.StatusNotFound))
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
					Day:         1,
					Month:       2,
					Year:        2022,
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
			Expect(updated.Day).To(Equal(savedExpense.Day))
			Expect(updated.Month).To(Equal(savedExpense.Month))
			Expect(updated.Year).To(Equal(savedExpense.Year))
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

				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
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
					Day:                1,
					Month:              1,
					Year:               2022,
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

				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(got).To(BeNil())
			})
		})
	})
})
