package auth_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/auth"
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
		accessTokenDuration            time.Duration
		ctx                            context.Context
		googleTokenValidatorMock       *auth.GoogleTokenAuth
	)

	BeforeEach(func() {
		userAuthenticableMock = &mocks.UserAuthenticable{}
		tokenizableMock = &mocks.Tokenizable{}
		googleTokenOpenIDValidatorMock = &mocks.OpenIDTokenValidable[validator.GoogleTokenClaims]{}
		accessTokenDuration = time.Duration(1 * time.Minute)
		ctx = context.Background()
		googleTokenValidatorMock = auth.NewGoogleTokenAuth(
			userAuthenticableMock,
			tokenizableMock,
			googleTokenOpenIDValidatorMock,
			accessTokenDuration,
		)
	})

	AfterEach(func() {
		T := GinkgoT()
		userAuthenticableMock.AssertExpectations(T)
		tokenizableMock.AssertExpectations(T)
	})

	Describe("Login", func() {
		It("takes google token to create token to login", func() {
			var (
				expectedToken      = faker.Paragraph()
				expectedUserEmail  = faker.Email()
				expectedUserName   = faker.Name()
				expectedUserAvatar = faker.URL()
				expectedUserID     = uuid.New()
				expectedFoundUser  = auth.UserData{
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
				expectedCreatedAccessToken   = faker.Paragraph()
				expectedAccessTokenExpiresAt = time.Now()
				expectedLoginOutput          = auth.LoginOutput{
					AccessToken:          expectedCreatedAccessToken,
					AccessTokenExpiresAt: expectedAccessTokenExpiresAt,
					User:                 &expectedFoundUser,
				}
			)
			googleTokenOpenIDValidatorMock.EXPECT().ValidateOpenIDToken(ctx, expectedToken).Return(&expectedGoogleClaims, nil)
			userAuthenticableMock.EXPECT().FindUserByEmail(ctx, expectedUserEmail).Return(&expectedFoundUser, nil)
			tokenizableMock.EXPECT().CreateAccessToken(&expectedAccessToken).Return(&auth.TokenInfo{
				Token:     expectedCreatedAccessToken,
				ExpiresAt: expectedAccessTokenExpiresAt,
			}, nil)

			got, err := googleTokenValidatorMock.Login(ctx, expectedToken)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(Equal(&expectedLoginOutput))
		})

		When("user is not register yet", func() {
			It("creates a new user to db and creates a new login token", func() {
				var (
					expectedToken      = faker.Paragraph()
					expectedUserEmail  = faker.Email()
					expectedUserName   = faker.Name()
					expectedUserAvatar = faker.URL()
					expectedUserID     = uuid.New()
					expectedUserToSave = auth.UserData{
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
					expectedCreatedAccessToken   = faker.Paragraph()
					expectedAccessTokenExpiresAt = time.Now()
					expectedLoginOutput          = auth.LoginOutput{
						AccessToken:          expectedCreatedAccessToken,
						AccessTokenExpiresAt: expectedAccessTokenExpiresAt,
						User: &auth.UserData{
							ID:     expectedUserID,
							Name:   expectedUserName,
							Email:  expectedUserEmail,
							Avatar: expectedUserAvatar,
						},
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

				got, err := googleTokenValidatorMock.Login(ctx, expectedToken)

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
