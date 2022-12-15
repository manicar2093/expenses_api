package tokens

import (
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

func (c *Paseto) CreateAccessToken(tokenDetails *auth.AccessToken) (string, error) {
	return c.createTokenWithClaims(tokenDetails.Expiration, map[string]interface{}{
		"user_id": tokenDetails.UserID,
	})
}

func (c *Paseto) CreateRefreshToken(tokenDetails *auth.RefreshToken) (string, error) {
	return c.createTokenWithClaims(tokenDetails.Expiration, map[string]interface{}{
		"session_id": tokenDetails.SessionID,
	})
}

func (c *Paseto) createTokenWithClaims(expiration time.Duration, claims map[string]interface{}) (string, error) {

	token := paseto.NewToken()
	now := time.Now()
	token.SetIssuedAt(now)
	token.SetExpiration(now.Add(expiration))
	for k, v := range claims {
		token.Set(k, v)
	}

	return token.V4Encrypt(c.symmetricKey, nil), nil
}

func (c *Paseto) ValidateToken(token string) error {
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

	return nil
}
