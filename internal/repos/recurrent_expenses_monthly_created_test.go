package repos_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/periodtypes"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("RecurrentExpense", func() {
	var (
		ctx                                            context.Context
		coll, expensesCountColl, recurrentExpensesColl *mongo.Collection
		repo                                           *repos.RecurrentExpensesMonthlyCreatedRepoImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		coll = conn.Collection("recurrent_expenses_monthly_created")
		expensesCountColl = conn.Collection(entities.ExpensesCountCollName)
		recurrentExpensesColl = conn.Collection(entities.RecurrentExpensesCollectonName)
		repo = repos.NewRecurrentExpensesCreatedMonthlyRepoImpl(conn)

	})

	Describe("Save", func() {
		It("storage a registry at db", func() {
			var (
				expectedRecurrentExpenseCreatedMonthly = entities.RecurrentExpensesMonthlyCreated{
					Month: 11,
					Year:  2022,
				}
			)
			err := repo.Save(ctx, &expectedRecurrentExpenseCreatedMonthly)
			var fromDB entities.RecurrentExpensesMonthlyCreated
			coll.FindOne(ctx, bson.D{{Key: "_id", Value: expectedRecurrentExpenseCreatedMonthly.ID}}).Decode(&fromDB) //nolint:errcheck

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedRecurrentExpenseCreatedMonthly.ID.IsZero()).To(BeFalse())

			testfunc.DeleteOneByObjectID(ctx, coll, expectedRecurrentExpenseCreatedMonthly.ID)
		})
	})

	Describe("FindByCurrentMonthAndYear", func() {

		It("returns found data", func() {
			var (
				expectedRecurrentExpenseLastCreationDate = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
				expectedRecurrentExpense                 = entities.RecurrentExpense{
					ID:               primitive.NewObjectID(),
					Name:             faker.Name(),
					Amount:           faker.Latitude(),
					Description:      faker.Paragraph(),
					Periodicity:      periodtypes.Daily,
					LastCreationDate: &expectedRecurrentExpenseLastCreationDate,
				}
				expectedMonth                          = uint(11)
				expectedYear                           = uint(2022)
				expectedRecurrentExpenseCreatedMonthly = entities.RecurrentExpensesMonthlyCreated{
					ID:    primitive.NewObjectID(),
					Month: expectedMonth,
					Year:  expectedYear,
				}
				expectedExpensesCount = entities.ExpensesCount{
					ID:                                primitive.NewObjectID(),
					RecurrentExpenseID:                expectedRecurrentExpense.ID,
					RecurrentExpense:                  &expectedRecurrentExpense,
					RecurrentExpensesMonthlyCreatedID: expectedRecurrentExpenseCreatedMonthly.ID,
					ExpensesRelatedIDs:                []primitive.ObjectID{primitive.NewObjectID()},
					TotalExpenses:                     1,
					TotalExpensesPaid:                 0,
				}
			)
			recurrentExpensesColl.InsertOne(ctx, &expectedRecurrentExpense) //nolint:errcheck
			expensesCountColl.InsertOne(ctx, &expectedExpensesCount)        //nolint:errcheck
			coll.InsertOne(ctx, &expectedRecurrentExpenseCreatedMonthly)    //nolint:errcheck

			got, err := repo.FindByCurrentMonthAndYear(ctx, expectedMonth, expectedYear)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.Month).To(Equal(expectedRecurrentExpenseCreatedMonthly.Month))
			Expect(got.Year).To(Equal(expectedRecurrentExpenseCreatedMonthly.Year))
			Expect(got.ExpensesCount[0].ID).To(Equal(expectedExpensesCount.ID))
			Expect(got.ExpensesCount[0].RecurrentExpense.Name).To(Equal(expectedRecurrentExpense.Name))

			testfunc.DeleteOneByObjectID(ctx, coll, expectedRecurrentExpenseCreatedMonthly.ID)
		})

		When("there is any data in db", func() {
			It("returns a NotFoundError", func() {
				got, err := repo.FindByCurrentMonthAndYear(ctx, 20, 1993)

				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(got).To(BeNil())

			})
		})
	})

	Describe("Update", func() {
		It("updates data in db with given instace", func() {
			var (
				expectedRecurrentExpenseCreatedMonthlySavedCreatedAt = time.Now()
				expectedRecurrentExpenseCreatedMonthlySavedID        = primitive.NewObjectID()
				expectedRecurrentExpenseCreatedMonthlySaved          = entities.RecurrentExpensesMonthlyCreated{
					ID:        expectedRecurrentExpenseCreatedMonthlySavedID,
					Month:     11,
					Year:      2022,
					CreatedAt: &expectedRecurrentExpenseCreatedMonthlySavedCreatedAt,
				}
				expectedUpdatedRecurrentExpenseCreatedMonthly = entities.RecurrentExpensesMonthlyCreated{
					ID:        expectedRecurrentExpenseCreatedMonthlySavedID,
					Month:     11,
					Year:      2022,
					CreatedAt: &expectedRecurrentExpenseCreatedMonthlySavedCreatedAt,
				}
			)
			inserted, _ := coll.InsertOne(ctx, expectedRecurrentExpenseCreatedMonthlySaved)

			err := repo.Update(ctx, &expectedUpdatedRecurrentExpenseCreatedMonthly)

			var fromDB entities.RecurrentExpensesMonthlyCreated
			coll.FindOne(ctx, bson.D{{Key: "_id", Value: inserted.InsertedID}}).Decode(&fromDB) //nolint:errcheck
			Expect(err).ToNot(HaveOccurred())
			Expect(fromDB.ID).To(Equal(expectedUpdatedRecurrentExpenseCreatedMonthly.ID))
			Expect(fromDB.Month).To(Equal(expectedUpdatedRecurrentExpenseCreatedMonthly.Month))
			Expect(fromDB.Year).To(Equal(expectedUpdatedRecurrentExpenseCreatedMonthly.Year))

			testfunc.DeleteOneByObjectID(ctx, coll, inserted.InsertedID.(primitive.ObjectID))
		})
	})

})
