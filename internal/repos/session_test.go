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
)

var _ = Describe("Session", func() {

	var (
		ctx            context.Context
		expectedUserID uuid.UUID
		expectedUser   entities.User
		repo           *repos.SessionsGormRepo
	)

	BeforeEach(func() {
		ctx = context.Background()
		expectedUserID = uuid.New()
		expectedUser = entities.User{
			ID:    expectedUserID,
			Email: faker.Email(),
		}
		conn.Create(&expectedUser)
		repo = repos.NewSessionsGormRepo(conn)
	})

	AfterEach(func() {
		conn.Delete(&expectedUser)
	})

	Describe("Save", func() {
		It("stores in db a instance of entities.Session", func() {
			var (
				sessionToSave = entities.Session{
					UserID:       expectedUserID,
					RefreshToken: faker.Paragraph(),
					UserAgent:    faker.Username(),
					ClientIP:     faker.IPv4(),
					ExpiresAt:    time.Now(),
				}
			)

			Expect(repo.Save(ctx, &sessionToSave)).To(Succeed())
		})
	})

	Describe("FindByID", func() {
		It("returns session with user data embeded", func() {
			var (
				expectedSessionID = uuid.New()
				sessionSaved      = entities.Session{
					ID:           expectedSessionID,
					UserID:       expectedUserID,
					RefreshToken: faker.Paragraph(),
					UserAgent:    faker.Username(),
					ClientIP:     faker.IPv4(),
					ExpiresAt:    time.Now(),
				}
			)
			conn.Create(&sessionSaved)
			defer conn.Delete(&sessionSaved)

			got, err := repo.FindByID(ctx, expectedSessionID)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.User.ID).To(Equal(expectedUserID))
		})
	})

})
