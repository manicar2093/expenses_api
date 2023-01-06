package sessions_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/sessions"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("Validate", func() {
	var (
		sessionFindableMock *mocks.SessionFindable
		ctx                 context.Context
		api                 *sessions.DefaultValidator
	)

	BeforeEach(func() {
		sessionFindableMock = mocks.NewSessionFindable(GinkgoT())
		api = sessions.NewDefaultValidator(sessionFindableMock)
	})

	Describe("ValidateSession", func() {
		It("check if a session is valid", func() {
			var (
				expectedSessionFound = entities.Session{
					ID:        uuid.New(),
					UserID:    uuid.New(),
					UserAgent: faker.Name(),
					ClientIP:  faker.IPv4(),
					CreatedAt: testfunc.ToPointer(time.Now()),
				}
				expectedValidationSessionInput = sessions.SessionValidationInput{
					SessionID:     expectedSessionFound.ID,
					FromUserAgent: expectedSessionFound.UserAgent,
					FromClientIP:  expectedSessionFound.ClientIP,
				}
			)
			sessionFindableMock.EXPECT().FindByID(ctx, expectedSessionFound.ID).Return(&expectedSessionFound, nil)

			got, err := api.ValidateSession(ctx, &expectedValidationSessionInput)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(Equal(&expectedSessionFound))
		})

		When("session is not found", func() {
			It("returns sessionFindable error", func() {
				var (
					expectedValidationSessionInput = sessions.SessionValidationInput{
						SessionID:     uuid.New(),
						FromUserAgent: faker.Name(),
						FromClientIP:  faker.IPv4(),
					}
					expectedSessionFindableError = &apperrors.NotFoundError{}
				)
				sessionFindableMock.EXPECT().FindByID(ctx, expectedValidationSessionInput.SessionID).Return(nil, expectedSessionFindableError)

				got, err := api.ValidateSession(ctx, &expectedValidationSessionInput)

				Expect(err).To(Equal(expectedSessionFindableError))
				Expect(got).To(BeNil())
			})
		})
	})

})
