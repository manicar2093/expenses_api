package repos_test

import (
	"context"
	"log"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"
)

var _ = Describe("RecurrentExpense", func() {

	var (
		ctx  context.Context
		repo *repos.RecurrentExpenseGormRepo
	)

	BeforeEach(func() {
		ctx = context.TODO()
		repo = repos.NewRecurrentExpenseGormRepo(conn)

	})

	Describe("Save", func() {
		It("saves an instance", func() {
			var (
				toSave = entities.RecurrentExpense{
					Name:   faker.Name(),
					Amount: faker.Latitude(),
					Description: null.StringFrom(
						faker.Paragraph(),
					),
				}
			)

			err := repo.Save(ctx, &toSave)
			defer conn.Delete(&toSave)

			log.Println(toSave.ID)

			Expect(err).ToNot(HaveOccurred())
			Expect(toSave.ID).ToNot(BeEmpty())
			Expect(toSave.CreatedAt).ToNot(BeZero())
			Expect(toSave.UpdatedAt).ToNot(BeZero())

		})

		When("recurrent expense exists", func() {
			It("returns an AlreadyExists", func() {
				var (
					expectedName = "testing"
					saved        = entities.RecurrentExpense{
						Name:   expectedName,
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					}
					toSave = entities.RecurrentExpense{
						Name:   expectedName,
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					}
				)
				conn.Create(&saved)
				defer conn.Delete(&saved)

				err := repo.Save(ctx, &toSave)
				defer conn.Delete(&toSave)

				Expect(err).To(BeAssignableToTypeOf(&repos.AlreadyExistsError{}))
				Expect(err.(*repos.AlreadyExistsError).Entity).To(Equal("Recurrent Expense"))
				Expect(err.(*repos.AlreadyExistsError).Identifier).To(Equal(saved.Name))
			})
		})
	})

	Describe("FindByName", func() {
		It("returns a pointer of found data", func() {
			var (
				expectedName = "testing"
				saved        = entities.RecurrentExpense{
					Name:   expectedName,
					Amount: faker.Latitude(),
					Description: null.StringFrom(
						faker.Paragraph(),
					),
				}
			)
			conn.Create(&saved)
			defer conn.Delete(&saved)

			got, err := repo.FindByName(ctx, expectedName)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.Name).To(Equal(expectedName))
		})
	})

	Describe("FindAll", func() {
		It("gets all registered recurrent expenses", func() {
			var (
				dataSaved = []*entities.RecurrentExpense{
					{
						Name:   faker.Name(),
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					},
					{
						Name:   faker.Name(),
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					},
					{
						Name:   faker.Name(),
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					},
				}
			)
			conn.Create(&dataSaved)
			defer conn.Delete(&dataSaved)

			got, err := repo.FindAll(ctx)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(HaveLen(len(dataSaved)))
		})
	})
})
