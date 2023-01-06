package auth_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/auth"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/validator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoogleAuth", func() {

	var (
		userAuthenticableMock          *mocks.UserAuthenticable
		tokenizableMock                *mocks.Tokenizable
		googleTokenOpenIDValidatorMock *mocks.OpenIDTokenValidable[validator.GoogleTokenClaims]
		sessionCreateableMock          *mocks.SessionCreateable
		accessTokenDuration            time.Duration
		ctx                            context.Context
		googleTokenValidatorMock       *auth.GoogleTokenAuth
	)

	BeforeEach(func() {
		T := GinkgoT()
		userAuthenticableMock = mocks.NewUserAuthenticable(T)
		tokenizableMock = mocks.NewTokenizable(T)
		googleTokenOpenIDValidatorMock = mocks.NewOpenIDTokenValidable[validator.GoogleTokenClaims](T)
		sessionCreateableMock = mocks.NewSessionCreateable(T)
		accessTokenDuration = time.Duration(1 * time.Minute)
		ctx = context.Background()
		googleTokenValidatorMock = auth.NewGoogleTokenAuth(
			userAuthenticableMock,
			tokenizableMock,
			googleTokenOpenIDValidatorMock,
			accessTokenDuration,
			sessionCreateableMock,
		)
	})

	Describe("Login", func() {

		var (
			expectedToken, expectedUserEmail, expectedUserName string
			expectedUserAvatar, expectedCreatedAccessToken     string
			expectedUserID                                     uuid.UUID
			expectedLoginInput                                 auth.LoginInput
			expectedSessionToCreate                            entities.Session
			expectedFoundUser                                  auth.UserData
			expectedGoogleClaims                               validator.GoogleTokenClaims
			expectedAccessToken                                auth.AccessToken
			expectedAccessTokenExpiresAt                       time.Time
			expectedLoginOutput                                auth.LoginOutput
		)

		BeforeEach(func() {
			expectedToken = faker.Paragraph()
			expectedUserEmail = faker.Email()
			expectedUserName = faker.Name()
			expectedUserAvatar = faker.URL()
			expectedUserID = uuid.New()
			expectedLoginInput = auth.LoginInput{
				Token:     expectedToken,
				UserAgent: faker.Name(),
				ClientIP:  faker.IPv4(),
			}
			expectedSessionToCreate = entities.Session{
				UserID:    expectedUserID,
				UserAgent: expectedLoginInput.UserAgent,
				ClientIP:  expectedLoginInput.ClientIP,
			}
			expectedFoundUser = auth.UserData{
				ID:     expectedUserID,
				Name:   expectedUserName,
				Email:  expectedUserEmail,
				Avatar: expectedUserAvatar,
			}
			expectedGoogleClaims = validator.GoogleTokenClaims{
				Email:         expectedUserEmail,
				EmailVerified: expectedUserEmail,
				Name:          expectedUserName,
				Picture:       expectedUserAvatar,
			}
			expectedAccessToken = auth.AccessToken{
				Expiration: accessTokenDuration,
				UserID:     expectedUserID,
			}
			expectedCreatedAccessToken = faker.Paragraph()
			expectedAccessTokenExpiresAt = time.Now()
			expectedLoginOutput = auth.LoginOutput{
				AccessToken:          expectedCreatedAccessToken,
				AccessTokenExpiresAt: expectedAccessTokenExpiresAt,
				RefreshToken:         expectedSessionToCreate.ID,
				User:                 &expectedFoundUser,
			}
		})

		It("takes google token to create token to login", func() {

			googleTokenOpenIDValidatorMock.EXPECT().ValidateOpenIDToken(ctx, expectedToken).Return(&expectedGoogleClaims, nil)
			userAuthenticableMock.EXPECT().FindUserByEmail(ctx, expectedUserEmail).Return(&expectedFoundUser, nil)
			tokenizableMock.EXPECT().CreateAccessToken(&expectedAccessToken).Return(&auth.TokenInfo{
				Token:     expectedCreatedAccessToken,
				ExpiresAt: expectedAccessTokenExpiresAt,
			}, nil)
			sessionCreateableMock.EXPECT().Create(ctx, &expectedSessionToCreate).Return(nil).Run(func(ctx context.Context, session *entities.Session) {
				expectedSessionToCreate.ID = uuid.New()
			})

			got, err := googleTokenValidatorMock.Login(ctx, &expectedLoginInput)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(Equal(&expectedLoginOutput))
		})

		When("user is not register yet", func() {
			It("creates a new user to db and creates a new login token", func() {
				var (
					expectedUserToSave = auth.UserData{
						Name:   expectedUserName,
						Email:  expectedUserEmail,
						Avatar: expectedUserAvatar,
					}
				)
				googleTokenOpenIDValidatorMock.EXPECT().ValidateOpenIDToken(ctx, expectedToken).Return(&expectedGoogleClaims, nil)
				userAuthenticableMock.EXPECT().FindUserByEmail(ctx, expectedUserEmail).Return(nil, &apperrors.NotFoundError{})
				userAuthenticableMock.EXPECT().CreateUser(ctx, &expectedUserToSave).Run(func(ctx context.Context, user *auth.UserData) {
					user.ID = expectedUserID
				}).Return(nil)
				tokenizableMock.EXPECT().CreateAccessToken(&expectedAccessToken).Return(&auth.TokenInfo{
					Token:     expectedCreatedAccessToken,
					ExpiresAt: expectedAccessTokenExpiresAt,
				}, nil)
				sessionCreateableMock.EXPECT().Create(ctx, &expectedSessionToCreate).Return(nil).Run(func(ctx context.Context, session *entities.Session) {
					expectedSessionToCreate.ID = uuid.New()
				})

				got, err := googleTokenValidatorMock.Login(ctx, &expectedLoginInput)

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal(&expectedLoginOutput))
			})
		})
	})

	Describe("RefreshToken", func() {
		It("creates a new access token from given refresh token", func() {

		})
	})

})
