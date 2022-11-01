package repos_test

import (
	"context"
	"time"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("RecurrentExpense", func() {
	var (
		ctx  context.Context
		coll *mongo.Collection
		repo *repos.RecurrentExpensesMonthlyCreatedRepoImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		coll = conn.Collection("recurrent_expenses_created_monthly")
		repo = repos.NewRecurrentExpensesCreatedMonthlyRepoImpl(conn)

	})

	Describe("Save", func() {
		It("storage a registry at db", func() {
			var (
				expectedCreatedAt                      = time.Now()
				expectedRecurrentExpenseCreatedMonthly = entities.RecurrentExpensesMonthlyCreated{
					Month: 11,
					Year:  2022,
					ExpensesCount: &entities.ExpensesCount{
						RecurrentExpenseID: primitive.NewObjectID(),
						ExpensesRelated: []primitive.ObjectID{
							primitive.NewObjectID(),
							primitive.NewObjectID(),
						},
						TotalExpenses:     2,
						TotalExpensesPaid: 0,
					},
					CreatedAt: &expectedCreatedAt,
				}
			)
			err := repo.Save(ctx, &expectedRecurrentExpenseCreatedMonthly)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedRecurrentExpenseCreatedMonthly.ID.IsZero()).To(BeFalse())

			testfunc.DeleteOneByObjectID(ctx, coll, expectedRecurrentExpenseCreatedMonthly.ID)
		})
	})
})
