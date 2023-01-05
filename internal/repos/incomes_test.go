package repos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
)

var _ = Describe("IncomesRepo", func() {
	var (
		ctx            context.Context
		repo           *repos.IncomesGormRepo
		expectedUserID uuid.UUID
		expectedUser   entities.User
	)

	BeforeEach(func() {
		ctx = context.TODO()
		repo = repos.NewIncomesGormRepo(conn)
		expectedUserID = uuid.New()
		expectedUser = entities.User{
			ID:    expectedUserID,
			Email: faker.Email(),
		}
		conn.Create(&expectedUser)
	})

	AfterEach(func() {
		conn.Delete(&expectedUser)
	})

	Describe("Save", func() {
		It("saves an entities.Income in database", func() {

			var (
				expectedName        = faker.Name()
				expectedAmount      = faker.Latitude()
				expectedDescription = faker.Sentence()
				expectedIncome      = entities.Income{
					UserID:      expectedUserID,
					Name:        expectedName,
					Amount:      expectedAmount,
					Day:         1,
					Month:       2,
					Year:        2022,
					Description: expectedDescription,
				}
			)

			err := repo.Save(ctx, &expectedIncome)
			defer conn.Delete(&expectedIncome)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedIncome.ID).ToNot(BeEmpty())
			Expect(expectedIncome.Day).ToNot(BeZero())
			Expect(expectedIncome.Month).ToNot(BeZero())
			Expect(expectedIncome.Year).ToNot(BeZero())
			Expect(expectedIncome.CreatedAt).ToNot(BeZero())
			Expect(expectedIncome.UpdatedAt).ToNot(BeZero())
		})
	})

})
