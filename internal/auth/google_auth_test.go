package auth_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/auth"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/sessions"
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
		sessionValidableMock           *mocks.SessionValidable
		userFindableMock               *mocks.UserFindable
		accessTokenDuration            time.Duration
		ctx                            context.Context
		googleTokenValidator           *auth.GoogleTokenAuth
	)

	BeforeEach(func() {
		T := GinkgoT()
		userAuthenticableMock = mocks.NewUserAuthenticable(T)
		tokenizableMock = mocks.NewTokenizable(T)
		googleTokenOpenIDValidatorMock = mocks.NewOpenIDTokenValidable[validator.GoogleTokenClaims](T)
		sessionCreateableMock = mocks.NewSessionCreateable(T)
		sessionValidableMock = mocks.NewSessionValidable(T)
		userFindableMock = mocks.NewUserFindable(T)
		accessTokenDuration = time.Duration(1 * time.Minute)
		ctx = context.Background()
		googleTokenValidator = auth.NewGoogleTokenAuth(
			userAuthenticableMock,
			tokenizableMock,
			googleTokenOpenIDValidatorMock,
			sessionCreateableMock,
			sessionValidableMock,
			userFindableMock,
			accessTokenDuration,
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

			got, err := googleTokenValidator.Login(ctx, &expectedLoginInput)

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

				got, err := googleTokenValidator.Login(ctx, &expectedLoginInput)

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal(&expectedLoginOutput))
			})
		})
	})

	Describe("RefreshToken", func() {
		It("creates a new access token from given refresh token", func() {
			var (
				expectedSessionUUID = uuid.New()
				expectedUserUUID    = uuid.New()
				expectedFoundUser   = auth.UserData{
					ID:     expectedUserUUID,
					Name:   faker.Name(),
					Email:  faker.Email(),
					Avatar: faker.URL(),
				}
				expectedValidateSessionCall = sessions.SessionValidationInput{
					SessionID:     expectedSessionUUID,
					FromUserAgent: faker.Name(),
					FromClientIP:  faker.IPv4(),
				}
				expectedSessionReturned = entities.Session{
					ID:        expectedSessionUUID,
					UserID:    expectedUserUUID,
					UserAgent: expectedValidateSessionCall.FromUserAgent,
					ClientIP:  expectedValidateSessionCall.FromClientIP,
				}
				expectedCreateAccessTokenCall = auth.AccessToken{
					UserID:     expectedFoundUser.ID,
					Expiration: accessTokenDuration,
				}
				expectedCreateAccessTokenReturn = auth.TokenInfo{
					Token:     faker.Paragraph(),
					ExpiresAt: time.Now(),
				}
				expectedRefreshTokenInput = auth.RefreshTokenInput{
					SessionID: expectedSessionUUID.String(),
					UserAgent: expectedValidateSessionCall.FromUserAgent,
					ClientIP:  expectedValidateSessionCall.FromClientIP,
				}
			)
			sessionValidableMock.EXPECT().ValidateSession(ctx, &expectedValidateSessionCall).Return(&expectedSessionReturned, nil)
			userFindableMock.EXPECT().FindUserByID(ctx, expectedFoundUser.ID).Return(&expectedFoundUser, nil)
			tokenizableMock.EXPECT().CreateAccessToken(&expectedCreateAccessTokenCall).Return(&expectedCreateAccessTokenReturn, nil)

			got, err := googleTokenValidator.RefreshToken(ctx, &expectedRefreshTokenInput)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.AccessToken).To(Equal(expectedCreateAccessTokenReturn.Token))
			Expect(got.AccessTokenExpiresAt).To(Equal(expectedCreateAccessTokenReturn.ExpiresAt))
			Expect(got.RefreshToken).To(Equal(expectedSessionUUID))

		})
	})

})
