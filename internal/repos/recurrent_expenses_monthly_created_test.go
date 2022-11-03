package repos_test

import (
	"context"
	"log"

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
		coll = conn.Collection("recurrent_expenses_monthly_created")
		repo = repos.NewRecurrentExpensesCreatedMonthlyRepoImpl(conn)

	})

	Describe("Save", func() {
		It("storage a registry at db", func() {
			var (
				expectedRecurrentExpenseCreatedMonthly = entities.RecurrentExpensesMonthlyCreated{
					Month: 11,
					Year:  2022,
					ExpensesCount: []*entities.ExpensesCount{
						{
							RecurrentExpenseID: primitive.NewObjectID(),
							ExpensesRelated: []primitive.ObjectID{
								primitive.NewObjectID(),
								primitive.NewObjectID(),
							},
							TotalExpenses:     2,
							TotalExpensesPaid: 0,
						},
					},
				}
			)
			err := repo.Save(ctx, &expectedRecurrentExpenseCreatedMonthly)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedRecurrentExpenseCreatedMonthly.ID.IsZero()).To(BeFalse())

			testfunc.DeleteOneByObjectID(ctx, coll, expectedRecurrentExpenseCreatedMonthly.ID)
		})
	})

	Describe("FindByMonthAndYear", func() {

		It("returns found data", func() {
			var (
				expectedMonth                          = uint(11)
				expectedYear                           = uint(2022)
				expectedRecurrentExpenseCreatedMonthly = entities.RecurrentExpensesMonthlyCreated{
					Month: expectedMonth,
					Year:  expectedYear,
					ExpensesCount: []*entities.ExpensesCount{
						{},
						{},
					},
				}
			)
			_, err := coll.InsertOne(ctx, &expectedRecurrentExpenseCreatedMonthly)
			if err != nil {
				log.Fatal(err)
			}

			got, err := repo.FindByMonthAndYear(ctx, expectedMonth, expectedYear)
			log.Println(*got)
			Expect(err).ToNot(HaveOccurred())
			Expect(got.Month).To(Equal(expectedRecurrentExpenseCreatedMonthly.Month))
			Expect(got.Year).To(Equal(expectedRecurrentExpenseCreatedMonthly.Year))

			testfunc.DeleteOneByObjectID(ctx, coll, expectedRecurrentExpenseCreatedMonthly.ID)
		})

		When("there is any data in db", func() {
			It("returns a NotFoundError", func() {
				got, err := repo.FindByMonthAndYear(ctx, 20, 1993)

				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
				Expect(got).To(BeNil())

			})
		})
	})
})
