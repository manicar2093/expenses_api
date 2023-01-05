package tokens

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/manicar2093/expenses_api/internal/auth"
)

var (
	ErrNotValidEncryptionKey = errors.New("encription key has no correct data to be used")
	ErrTokenExpired          = errors.New("token has expired")
)

type Paseto struct {
	symmetricKey paseto.V4SymmetricKey
}

func NewPaseto(symmetricKey string) *Paseto {
	encryptionKey, err := paseto.V4SymmetricKeyFromHex(symmetricKey)
	if err != nil {
		panic(ErrNotValidEncryptionKey)
	}
	return &Paseto{
		symmetricKey: encryptionKey,
	}
}

func (c *Paseto) CreateAccessToken(tokenDetails *auth.AccessToken) (*auth.TokenInfo, error) {
	token, expiresAt, err := c.createTokenWithClaims(tokenDetails.Expiration, map[string]interface{}{
		"user_id": tokenDetails.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &auth.TokenInfo{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (c *Paseto) CreateRefreshToken(tokenDetails *auth.RefreshToken) (*auth.TokenInfo, error) {
	token, expiresAt, err := c.createTokenWithClaims(tokenDetails.Expiration, map[string]interface{}{
		"session_id": tokenDetails.SessionID,
	})
	if err != nil {
		return nil, err
	}
	return &auth.TokenInfo{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (c *Paseto) createTokenWithClaims(expiration time.Duration, claims map[string]interface{}) (string, time.Time, error) {
	var (
		token     = paseto.NewToken()
		now       = time.Now()
		expiresAt = now.Add(expiration)
	)
	token.SetIssuedAt(now)
	token.SetExpiration(expiresAt)
	for k, v := range claims {
		if err := token.Set(k, v); err != nil {
			return "", time.Time{}, err
		}
	}

	return token.V4Encrypt(c.symmetricKey, nil), expiresAt, nil
}

func (c *Paseto) ValidateToken(ctx context.Context, token string, output interface{}) error {
	parser := paseto.NewParser()

	validatedToken, err := parser.ParseV4Local(c.symmetricKey, token, nil)
	if err != nil {
		return err
	}

	expirateAt, err := validatedToken.GetExpiration()
	if err != nil {
		return err
	}

	if time.Now().After(expirateAt) {
		return ErrTokenExpired
	}

	if err := json.Unmarshal(validatedToken.ClaimsJSON(), &output); err != nil {
		return err
	}

	return nil
}
