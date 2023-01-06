package repos_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("Sessions", func() {

	var (
		ctx  context.Context
		repo *repos.SessionsGormRepo
	)

	BeforeEach(func() {
		repo = repos.NewSessionGormRepo(conn)
	})

	Describe("FindByID", func() {
		It("retreives session found with given id", func() {
			var (
				userStored = entities.User{
					ID:    uuid.New(),
					Email: faker.Email(),
				}
				sessionStored = entities.Session{
					ID:        uuid.New(),
					UserID:    userStored.ID,
					UserAgent: faker.Name(),
					ClientIP:  faker.IPv4(),
					CreatedAt: testfunc.ToPointer(time.Now()),
				}
			)
			conn.Create(&userStored)
			conn.Create(&sessionStored)
			defer conn.Delete(&userStored)
			defer conn.Delete(&sessionStored)

			got, err := repo.FindByID(ctx, sessionStored.ID)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(BeAssignableToTypeOf(&entities.Session{}))
			Expect(got.ID).To(Equal(sessionStored.ID))
		})

		When("session is not found", func() {
			It("returns a apperrors.NotFoundError", func() {
				got, err := repo.FindByID(ctx, uuid.New())

				Expect(err).To(BeAssignableToTypeOf(&apperrors.NotFoundError{}))
				Expect(got).To(BeNil())
			})
		})
	})

})
