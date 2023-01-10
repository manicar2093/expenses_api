package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/sessions"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

type (
	GoogleTokenAuth struct {
		SessionCreateable
		sessions.SessionValidable
		UserFindable
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
	sessionCreateable SessionCreateable,
	sessionValidable sessions.SessionValidable,
	userFindable UserFindable,
	accessTokenDuration time.Duration,
) *GoogleTokenAuth {
	return &GoogleTokenAuth{
		SessionCreateable:    sessionCreateable,
		userAuthenticable:    userAuthenticable,
		tokenizable:          tokenizable,
		openIDTokenValidable: openIDTokenValidable,
		SessionValidable:     sessionValidable,
		UserFindable:         userFindable,
		accessTokenDuration:  accessTokenDuration,
	}
}

func (c *GoogleTokenAuth) Login(ctx context.Context, loginInput *LoginInput) (*LoginOutput, error) {
	googleClaims, err := c.openIDTokenValidable.ValidateOpenIDToken(ctx, loginInput.Token)
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
	sessionToCreate := entities.Session{
		UserID:    userFound.ID,
		UserAgent: loginInput.UserAgent,
		ClientIP:  loginInput.ClientIP,
	}
	if err := c.SessionCreateable.Create(ctx, &sessionToCreate); err != nil {
		return nil, err
	}
	return &LoginOutput{
		AccessToken:          tokenInfo.Token,
		AccessTokenExpiresAt: tokenInfo.ExpiresAt,
		RefreshToken:         sessionToCreate.ID,
		User:                 userFound,
	}, nil
}

func (c *GoogleTokenAuth) RefreshToken(ctx context.Context, refreshTokenInput *RefreshTokenInput) (*LoginOutput, error) {
	sessionIDAsUUID := uuid.MustParse(refreshTokenInput.SessionID)
	sessionInfo, err := c.ValidateSession(ctx, &sessions.SessionValidationInput{
		SessionID:     sessionIDAsUUID,
		FromUserAgent: refreshTokenInput.UserAgent,
		FromClientIP:  refreshTokenInput.ClientIP,
	})
	if err != nil {
		return nil, err
	}
	userFound, err := c.FindUserByID(ctx, sessionInfo.UserID)
	if err != nil {
		return nil, err
	}
	tokenInfo, err := c.tokenizable.CreateAccessToken(&AccessToken{UserID: userFound.ID, Expiration: c.accessTokenDuration})
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:          tokenInfo.Token,
		AccessTokenExpiresAt: tokenInfo.ExpiresAt,
		RefreshToken:         sessionIDAsUUID,
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
