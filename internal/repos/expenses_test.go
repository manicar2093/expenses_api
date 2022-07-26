package repos_test

import (
	"context"
	"log"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/connections"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("ExpensesImpl", func() {
	var (
		ctx            context.Context
		conn           *mongo.Database
		timeGetterMock *mocks.TimeGetable
		repo           *repos.ExpensesRepositoryImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		conn = connections.GetMongoConn()
		timeGetterMock = &mocks.TimeGetable{}
		repo = repos.NewExpensesRepositoryImpl(conn, timeGetterMock)

	})
	AfterEach(func() {
		conn.Drop(ctx)
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

			var expenseRegistered entities.Expense
			conn.Collection("expenses").FindOne(ctx, bson.M{"_id": expectedExpense.ID}).Decode(&expenseRegistered)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedExpense.ID).To(BeAssignableToTypeOf(primitive.ObjectID{}))
			Expect(expectedExpense.Day).ToNot(BeZero())
			Expect(expectedExpense.Month).ToNot(BeZero())
			Expect(expectedExpense.Year).ToNot(BeZero())
			Expect(expectedExpense.CreatedAt).ToNot(BeZero())
			Expect(expectedExpense.CreatedAt).To(Equal(expenseRegistered.CreatedAt))
			Expect(expectedExpense.UpdatedAt).To(BeNil())

		})
	})

	Describe("GetCurrentMonthExpenses", func() {
		It("returns all expenses by current month", func() {
			expectedTimeReturn := time.Date(2022, time.July, 1, 0, 0, 0, 0, time.Local)
			log.Println(expectedTimeReturn)
			timeGetterMock.EXPECT().GetCurrentTime().Return(expectedTimeReturn)
			expenses_created := []interface{}{
				bson.D{{Key: "month", Value: uint(time.July)}},
				bson.D{{Key: "month", Value: uint(time.July)}},
				bson.D{{Key: "month", Value: uint(time.July)}},
				bson.D{{Key: "month", Value: uint(time.March)}},
			}
			conn.Collection("expenses").InsertMany(ctx, expenses_created)
			got, err := repo.GetCurrentMonthExpenses(ctx)

			Expect(err).ToNot(HaveOccurred())
			Expect(*got).To(HaveLen(3))
		})
	})
})
