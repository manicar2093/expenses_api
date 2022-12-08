package repos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"
)

var _ = Describe("User", func() {

	var (
		ctx  context.Context
		repo *repos.UserGormRepo
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = repos.NewUserGormRepo(conn)
	})

	Describe("Save", func() {
		It("store in db a new user", func() {
			var (
				expectedUserToSave = entities.User{
					Name:     null.StringFrom(faker.Name()),
					Lastname: null.StringFrom(faker.LastName()),
					Email:    faker.Email(),
					Avatar:   null.StringFrom(faker.URL()),
				}
			)
			defer conn.Delete(&expectedUserToSave)

			err := repo.Save(ctx, &expectedUserToSave)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedUserToSave.ID.String()).ToNot(BeEmpty())
		})
	})

	Describe("FindUserByEmail", func() {
		It("return entities.User when is registried", func() {
			var (
				expectedEmail     = faker.Email()
				expectedUserSaved = entities.User{
					Name:     null.StringFrom(faker.Name()),
					Lastname: null.StringFrom(faker.LastName()),
					Email:    expectedEmail,
					Avatar:   null.StringFrom(faker.URL()),
				}
			)
			conn.Create(&expectedUserSaved)
			defer conn.Delete(&expectedUserSaved)

			got, err := repo.FindUserByEmail(ctx, expectedEmail)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(BeAssignableToTypeOf(&entities.User{}))
			Expect(got.Email).To(Equal(expectedEmail))
		})

		When("user is not found", func() {
			It("return a repos.NotFoundErr", func() {
				var (
					expectedEmail = faker.Email()
				)

				got, err := repo.FindUserByEmail(ctx, expectedEmail)

				Expect(got).To(BeNil())
				Expect(err).To(BeAssignableToTypeOf(&repos.NotFoundError{}))
			})
		})
	})

})
