package repos_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
)

var _ = Describe("ExpensesImpl", func() {
	var (
		ctx  context.Context
		repo *repos.ExpensesRepositoryImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		repo = repos.NewExpensesRepositoryImpl(conn)

	})

	Describe("Save", func() {
		It("saves an entities.Expense in database", func() {

			var (
				expectedName        = faker.Name()
				expectedAmount      = faker.Latitude()
				expectedDescription = faker.Sentence()
				expectedExpense     = entities.Expense{
					Name:        expectedName,
					Amount:      expectedAmount,
					Description: expectedDescription,
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
			Expect(expectedExpense.CreatedAt).ToNot(BeZero())
			Expect(expectedExpense.UpdatedAt).To(BeNil())

			testfunc.DeleteOneByObjectID(ctx, conn.Collection("expenses"), expectedExpense.ID)
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
			inserted, _ := conn.Collection("expenses").InsertMany(ctx, expensesCreated)
			got, err := repo.GetExpensesByMonth(ctx, time.July)

			Expect(err).ToNot(HaveOccurred())
			Expect(*got).To(HaveLen(3))

			testfunc.DeleteManyByObjectID(ctx, conn.Collection("expenses"), inserted)
		})
	})
})
