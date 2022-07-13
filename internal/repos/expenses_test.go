package repos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/connections"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("ExpensesImpl", func() {
	var (
		ctx  context.Context
		conn *mongo.Database
		repo *repos.ExpensesRepositoryImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		conn = connections.GetMongoConn()
		repo = repos.NewExpensesRepositoryImpl(conn)

	})
	AfterEach(func() {
		conn.Drop(ctx)
	})

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
		Expect(expectedExpense.CreatedAt).ToNot(BeZero())
		Expect(expectedExpense.CreatedAt).To(Equal(expenseRegistered.CreatedAt))
		Expect(expectedExpense.UpdatedAt).To(BeNil())

	})
})
