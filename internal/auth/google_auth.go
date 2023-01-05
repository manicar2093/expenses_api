package auth

import (
	"context"
	"errors"
	"time"

	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

type (
	GoogleTokenAuth struct {
		userAuthenticable    UserAuthenticable
		tokenizable          Tokenizable
		openIDTokenValidable OpenIDTokenValidable[validator.GoogleTokenClaims]
		accessTokenDuration  time.Duration
	}
)

func NewGoogleTokenAuth(
	userAuthenticable UserAuthenticable,
	tokenizable Tokenizable,
	openIDTokenValidable OpenIDTokenValidable[validator.GoogleTokenClaims],
	accessTokenDuration time.Duration,
) *GoogleTokenAuth {
	return &GoogleTokenAuth{
		userAuthenticable:    userAuthenticable,
		tokenizable:          tokenizable,
		openIDTokenValidable: openIDTokenValidable,
		accessTokenDuration:  accessTokenDuration,
	}
}

func (c *GoogleTokenAuth) Login(ctx context.Context, token string) (*LoginOutput, error) {
	googleClaims, err := c.openIDTokenValidable.ValidateOpenIDToken(ctx, token)
	if err != nil {
		return nil, err
	}
	userFound, err := c.findUserOrCreate(ctx, googleClaims)
	if err != nil {
		return nil, err
	}
	tokenInfo, err := c.tokenizable.CreateAccessToken(&AccessToken{
		UserID:     userFound.ID,
		Expiration: c.accessTokenDuration,
	})
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:          tokenInfo.Token,
		AccessTokenExpiresAt: tokenInfo.ExpiresAt,
		User:                 userFound,
	}, nil
}

func (c *GoogleTokenAuth) findUserOrCreate(ctx context.Context, googleClaims *validator.GoogleTokenClaims) (*UserData, error) {
	userFound, err := c.userAuthenticable.FindUserByEmail(ctx, googleClaims.Email)
	if err != nil {
		var notFoundError *apperrors.NotFoundError
		if errors.As(err, &notFoundError) {
			userToSave := &UserData{
				Name:   googleClaims.Name,
				Email:  googleClaims.Email,
				Avatar: googleClaims.Picture,
			}
			if err := c.userAuthenticable.CreateUser(ctx, userToSave); err != nil {
				return nil, err
			}
			return userToSave, nil
		}
		return nil, err
	}
	return userFound, nil
}
