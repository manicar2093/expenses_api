package repos_test

import (
	"context"
	"log"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/period"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"
)

var _ = Describe("RecurrentExpense", func() {

	var (
		ctx            context.Context
		repo           *repos.RecurrentExpenseGormRepo
		expectedUserID uuid.UUID
		expectedUser   entities.User
	)

	BeforeEach(func() {
		ctx = context.TODO()
		repo = repos.NewRecurrentExpenseGormRepo(conn)
		expectedUserID = uuid.New()
		expectedUser = entities.User{
			ID:       expectedUserID,
			Name:     null.NewString(faker.Name(), true),
			Lastname: null.NewString(faker.LastName(), true),
			Email:    faker.Email(),
		}
		conn.Create(&expectedUser)
	})

	AfterEach(func() {
		conn.Delete(&expectedUser)
	})

	Describe("Create", func() {
		It("saves an instance", func() {
			var (
				toSave = entities.RecurrentExpense{
					UserID: expectedUserID,
					Name:   faker.Name(),
					Amount: faker.Latitude(),
					Description: null.StringFrom(
						faker.Paragraph(),
					),
					Periodicity: period.Daily,
				}
			)
			defer conn.Delete(&toSave)

			err := repo.Create(ctx, &toSave)

			log.Println(toSave.ID)

			Expect(err).ToNot(HaveOccurred())
			Expect(toSave.ID).ToNot(BeEmpty())
			Expect(toSave.CreatedAt).ToNot(BeZero())
			Expect(toSave.UpdatedAt).To(BeZero())

		})

		When("recurrent expense exists", func() {
			It("returns an AlreadyExists", func() {
				var (
					expectedName = "testing"
					saved        = entities.RecurrentExpense{
						UserID: expectedUserID,
						Name:   expectedName,
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					}
					toSave = entities.RecurrentExpense{
						UserID: expectedUserID,
						Name:   expectedName,
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					}
				)
				conn.Create(&saved)
				defer conn.Delete(&saved)

				err := repo.Create(ctx, &toSave)
				defer conn.Delete(&toSave)

				Expect(err).To(BeAssignableToTypeOf(&apperrors.AlreadyExistsError{}))
				Expect(err.(*apperrors.AlreadyExistsError).Entity).To(Equal("Recurrent Expense"))
				Expect(err.(*apperrors.AlreadyExistsError).Identifier).To(Equal(saved.Name))
			})
		})
	})

	Describe("FindByName", func() {
		It("returns a pointer of found data", func() {
			var (
				expectedName = "testing"
				saved        = entities.RecurrentExpense{
					UserID: expectedUserID,
					Name:   expectedName,
					Amount: faker.Latitude(),
					Description: null.StringFrom(
						faker.Paragraph(),
					),
				}
			)
			conn.Create(&saved)
			defer conn.Delete(&saved)

			got, err := repo.FindByName(ctx, expectedName, expectedUserID)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.Name).To(Equal(expectedName))
		})
	})

	Describe("FindAll", func() {
		It("gets all registered recurrent expenses", func() {
			var (
				dataSaved = []*entities.RecurrentExpense{
					{
						UserID: expectedUserID,
						Name:   faker.Name(),
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					},
					{
						UserID: expectedUserID,
						Name:   faker.Name(),
						Amount: faker.Latitude(),
						Description: null.StringFrom(
							faker.Paragraph(),
						),
					},
					{
						UserID: expectedUserID,
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

			got, err := repo.FindAll(ctx, expectedUserID)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(HaveLen(len(dataSaved)))
		})
	})
})
