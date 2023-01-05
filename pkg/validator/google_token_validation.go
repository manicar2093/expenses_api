package validator

import (
	"context"
	"fmt"
	"net/http"

	gohttpclient "github.com/bozd4g/go-http-client"
)

const googleValidationURL = "https://oauth2.googleapis.com/tokeninfo?id_token="

var (
	ErrGoogleLogin = fmt.Errorf("google does not validate this token")
)

type (
	GoogleTokenClaims struct {
		Iss           string `json:"iss"`
		Azp           string `json:"azp"`
		Aud           string `json:"aud"`
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Locale        string `json:"locale"`
		Iat           string `json:"iat"`
		Exp           string `json:"exp"`
		Alg           string `json:"alg"`
		Kid           string `json:"kid"`
		Typ           string `json:"typ"`
	}
	GoogleTokenValidator[T GoogleTokenClaims] struct{}
)

func NewGoogleTokenValidator() *GoogleTokenValidator[GoogleTokenClaims] {
	return &GoogleTokenValidator[GoogleTokenClaims]{}
}

func (c *GoogleTokenValidator[GoogleTokenClaims]) ValidateOpenIDToken(ctx context.Context, token string) (*GoogleTokenClaims, error) {
	client := gohttpclient.New(fmt.Sprintf("%s%s", googleValidationURL, token))
	res, err := client.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	if res.Status() != http.StatusOK {
		return nil, ErrGoogleLogin
	}

	var tokenClaims GoogleTokenClaims
	if err := res.Unmarshal(&tokenClaims); err != nil {
		return nil, err
	}

	return &tokenClaims, nil
}
