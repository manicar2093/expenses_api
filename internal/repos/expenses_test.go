package repos_test

import (
	"context"
	"net/http"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/converters"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("ExpensesImpl", func() {
	var (
		ctx  context.Context
		coll *mongo.Collection
		repo *repos.ExpensesRepositoryImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		coll = conn.Collection("expenses")
		repo = repos.NewExpensesRepositoryImpl(conn)

	})

	Describe("Save", func() {
		It("saves an entities.Expense in database", func() {

			var (
				expectedRecurrentExpenseID = primitive.NewObjectID()
				expectedName               = faker.Name()
				expectedAmount             = faker.Latitude()
				expectedDescription        = faker.Sentence()
				expectedExpense            = entities.Expense{
					Name:               expectedName,
					Amount:             expectedAmount,
					Description:        expectedDescription,
					RecurrentExpenseID: &expectedRecurrentExpenseID,
				}
			)

			err := repo.Save(ctx, &expectedExpense)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedExpense.ID.String()).ToNot(BeEmpty())
			Expect(expectedExpense.Day).ToNot(BeZero())
			Expect(expectedExpense.Month).ToNot(BeZero())
			Expect(expectedExpense.Year).ToNot(BeZero())
			Expect(expectedExpense.IsPaid).To(BeFalse())
			Expect(expectedExpense.IsRecurrent).To(BeFalse())
			Expect(expectedExpense.RecurrentExpenseID).To(Equal(&expectedRecurrentExpenseID))
			Expect(expectedExpense.CreatedAt).ToNot(BeZero())
			Expect(expectedExpense.UpdatedAt).To(BeNil())

			// testfunc.DeleteOneByObjectID(ctx, coll, expectedExpense.ID)
		})

		When("createdAt is given set", func() {
			It("respects given time", func() {
				var (
					expectedName        = faker.Name()
					expectedAmount      = faker.Latitude()
					expectedDescription = faker.Sentence()
					expectedCreatedAt   = time.Date(2022, time.August, 0, 0, 0, 0, 0, time.Local)
					expectedExpense     = entities.Expense{
						Name:        expectedName,
						Amount:      expectedAmount,
						Description: expectedDescription,
						CreatedAt:   &expectedCreatedAt,
					}
				)

				err := repo.Save(ctx, &expectedExpense)

				Expect(err).ToNot(HaveOccurred())
				Expect(expectedExpense.ID.String()).ToNot(BeEmpty())
				Expect(expectedExpense.Day).ToNot(BeZero())
				Expect(expectedExpense.Month).ToNot(BeZero())
				Expect(expectedExpense.Year).ToNot(BeZero())
				Expect(expectedExpense.IsPaid).To(BeFalse())
				Expect(expectedExpense.IsRecurrent).To(BeFalse())
				Expect(expectedExpense.CreatedAt).To(Equal(&expectedCreatedAt))
				Expect(expectedExpense.UpdatedAt).To(BeNil())

				testfunc.DeleteOneByObjectID(ctx, coll, expectedExpense.ID)
			})
		})
	})

	Describe("GetExpensesByMonth", func() {
		It("returns all expenses by current month", func() {
			expensesCreated := []interface{}{
				bson.D{{Key: "month", Value: uint(time.July)}},
				bson.D{{Key: "month", Value: uint(time.July)}},
				bson.D{{Key: "month", Value: uint(time.July)}},
				bson.D{{Key: "month", Value: uint(time.March)}},
			}
			inserted, _ := coll.InsertMany(ctx, expensesCreated)
			got, err := repo.GetExpensesByMonth(ctx, time.July)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(HaveLen(3))
			Expect(got[0]).To(BeAssignableToTypeOf(&entities.Expense{}))

			testfunc.DeleteManyByObjectID(ctx, coll, inserted)
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
				expectedName   = faker.Name()
				expectedStatus = true
				mockData       = entities.Expense{
					Name:        expectedName,
					IsRecurrent: true,
					IsPaid:      expectedStatus,
				}
			)
			inserted, _ := coll.InsertOne(ctx, mockData)
			expectedID := inserted.InsertedID.(primitive.ObjectID)

			err := repo.UpdateIsPaidByExpenseID(ctx, expectedID.Hex(), expectedStatus)

			var changed entities.Expense
			coll.FindOne(ctx, bson.D{{Key: "_id", Value: expectedID}}).Decode(&changed)
			Expect(err).ToNot(HaveOccurred())
			Expect(changed.Name).To(Equal(expectedName))
			Expect(changed.IsPaid).To(Equal(expectedStatus))

			testfunc.DeleteOneByObjectID(ctx, coll, expectedID)
		})

		When("expense is not recurrent", func() {
			It("returns a NotFoundError", func() {
				var (
					expectedName   = faker.Name()
					expectedStatus = true
					mockData       = entities.Expense{
						Name:   expectedName,
						IsPaid: expectedStatus,
					}
				)
				inserted, _ := coll.InsertOne(ctx, mockData)
				expectedID := inserted.InsertedID.(primitive.ObjectID)

				err := repo.UpdateIsPaidByExpenseID(ctx, expectedID.Hex(), expectedStatus)

				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(err.(*repos.NotFoundError).StatusCode()).To(Equal(http.StatusNotFound))

				testfunc.DeleteOneByObjectID(ctx, coll, expectedID)
			})
		})

		When("expenseID is not primitive.ObjectID", func() {
			It("returns an error", func() {
				var (
					expectedID = "not_an_object_id"
				)

				err := repo.UpdateIsPaidByExpenseID(ctx, expectedID, false)

				Expect(err).To(BeAssignableToTypeOf(&converters.IDNotValidIDError{}))
			})
		})
	})

	Describe("FindByNameAndMonthAndIsRecurrent", func() {
		It("finds a expense that is recurrent by its name", func() {
			var (
				expectedExpenseName = faker.Name()
				expectedMockData    = []interface{}{
					entities.Expense{Name: expectedExpenseName, IsRecurrent: true, Month: 8},
					entities.Expense{Name: faker.Name(), IsRecurrent: true},
					entities.Expense{Name: faker.Name(), IsRecurrent: false},
				}
			)
			inserted, _ := coll.InsertMany(ctx, expectedMockData)
			defer testfunc.DeleteManyByObjectID(ctx, coll, inserted)

			got, err := repo.FindByNameAndMonthAndIsRecurrent(ctx, 8, expectedExpenseName)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.Name).To(Equal(expectedExpenseName))
			Expect(got.IsRecurrent).To(BeTrue())
		})

		When("expense does not exist", func() {
			It("return an repos.NotFoundError", func() {
				var (
					expectedExpenseName = faker.Name()
					expectedMockData    = []interface{}{
						entities.Expense{Name: expectedExpenseName, IsRecurrent: true},
						entities.Expense{Name: faker.Name(), IsRecurrent: true},
						entities.Expense{Name: faker.Name(), IsRecurrent: false},
					}
				)
				inserted, _ := coll.InsertMany(ctx, expectedMockData)
				defer testfunc.DeleteManyByObjectID(ctx, coll, inserted)

				got, err := repo.FindByNameAndMonthAndIsRecurrent(ctx, 8, expectedExpenseName)

				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(got).To(BeNil())
			})
		})
	})

	Describe("GetExpenseStatusByID", func() {
		It("finds a expense by ID retriving just is_paid and its ID", func() {
			var (
				expectedStatus = true
				expectedID     = primitive.NewObjectID()
				mockData       = entities.Expense{
					ID:     expectedID,
					Name:   faker.Name(),
					IsPaid: expectedStatus,
				}
			)
			inserted, _ := coll.InsertOne(ctx, mockData)

			got, err := repo.GetExpenseStatusByID(ctx, expectedID.Hex())

			Expect(err).ToNot(HaveOccurred())
			Expect(got.ID).To(Equal(inserted.InsertedID))
			Expect(got.IsPaid).To(Equal(expectedStatus))

			testfunc.DeleteOneByObjectID(ctx, coll, expectedID)
		})

		When("expense is not found", func() {
			It("returns a NotFound error", func() {
				var (
					expectedID = primitive.NewObjectID().Hex()
				)

				got, err := repo.GetExpenseStatusByID(ctx, expectedID)

				Expect(got).To(BeNil())
				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(err.(*repos.NotFoundError).StatusCode()).To(Equal(http.StatusNotFound))

			})
		})

		When("expenseID is not primitive.ObjectID", func() {
			It("returns an error", func() {
				var (
					expectedID = "not_an_object_id"
				)

				got, err := repo.GetExpenseStatusByID(ctx, expectedID)

				Expect(got).To(BeNil())
				Expect(err).To(BeAssignableToTypeOf(&converters.IDNotValidIDError{}))
			})
		})
	})
})
