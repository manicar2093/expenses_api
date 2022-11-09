package repos_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("ExpensesCount", func() {

	var (
		coll *mongo.Collection
		repo *repos.ExpensesCountMongoRepo
		ctx  context.Context
	)

	BeforeEach(func() {
		coll = conn.Collection(entities.ExpensesCountCollName)
		ctx = context.Background()
		repo = repos.NewExpensesCountMongoRepo(conn)
	})

	Describe("Save", func() {
		It("saves a instance of ExpensesCount", func() {
			var (
				expectedRecurrentExpenseID                = primitive.NewObjectID()
				expectedRecurrentExpensesMonthlyCreatedID = primitive.NewObjectID()
				expectedExpensesCount                     = entities.ExpensesCount{
					RecurrentExpenseID:                expectedRecurrentExpenseID,
					RecurrentExpensesMonthlyCreatedID: expectedRecurrentExpensesMonthlyCreatedID,
					ExpensesRelatedIDs:                []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
					TotalExpenses:                     2,
					TotalExpensesPaid:                 0,
				}
			)

			err := repo.Save(ctx, &expectedExpensesCount)
			var inDB entities.ExpensesCount
			coll.FindOne(ctx, bson.D{{Key: "_id", Value: expectedExpensesCount.ID}}).Decode(&inDB) //nolint:errcheck

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedExpensesCount.ID.IsZero()).To(BeFalse())
			Expect(inDB).To(Equal(expectedExpensesCount))

			testfunc.DeleteOneByObjectID(ctx, coll, expectedExpensesCount.ID)
		})
	})

})
