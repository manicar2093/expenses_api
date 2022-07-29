package repos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("RecurrentExpense", func() {
	var (
		ctx  context.Context
		repo *repos.RecurrentExpenseRepoImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		repo = repos.NewRecurrentExpenseRepoImpl(conn)

	})

	It("saves an instance", func() {
		var (
			toSave = entities.RecurrentExpense{
				Name:        faker.Name(),
				Amount:      faker.Latitude(),
				Description: faker.Paragraph(),
			}
		)

		err := repo.Save(ctx, &toSave)

		Expect(err).ToNot(HaveOccurred())
		Expect(toSave.ID.String()).ToNot(BeEmpty())
		Expect(toSave.CreatedAt).ToNot(BeZero())
		Expect(toSave.UpdatedAt).To(BeNil())

		testfunc.DeleteOneByObjectID(ctx, conn.Collection("recurrent_expenses"), toSave.ID)
	})

	When("recurrent expense exists", func() {
		It("returns an AlreadyExists", func() {
			var (
				expectedName = "testing"
				Saved        = entities.RecurrentExpense{
					Name:        expectedName,
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
				}
				toSave = entities.RecurrentExpense{
					Name:        expectedName,
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
				}
			)
			repo.Save(ctx, &Saved)

			err := repo.Save(ctx, &toSave)

			Expect(err).To(BeAssignableToTypeOf(&repos.AlreadyExistsError{}))
			Expect(err.(*repos.AlreadyExistsError).Entity).To(Equal("RecurrentExpense"))
			Expect(err.(*repos.AlreadyExistsError).Identifier).To(Equal(Saved.Name))

			coll := conn.Collection("recurrent_expenses")
			testfunc.DeleteOneByObjectID(ctx, coll, Saved.ID)
			testfunc.DeleteOneByObjectID(ctx, coll, toSave.ID)
		})
	})

})
