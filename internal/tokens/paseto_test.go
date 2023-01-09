package tokens_test

import (
	"context"
	"log"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/auth"
	"github.com/manicar2093/expenses_api/internal/tokens"
)

var _ = Describe("Paseto", func() {

	var (
		symmetricKey string
		api          *tokens.Paseto
	)

	BeforeEach(func() {
		symmetricKey = paseto.NewV4SymmetricKey().ExportHex()
		api = tokens.NewPaseto(symmetricKey)
	})

	Describe("CreateAccessToken", func() {
		It("creates an access token", func() {
			var (
				tokenDetails = auth.AccessToken{
					Expiration: time.Duration(1 * time.Second),
					UserID:     uuid.New(),
				}
			)

			token, err := api.CreateAccessToken(&tokenDetails)

			Expect(err).ToNot(HaveOccurred())
			Expect(token.Token).ToNot(BeEmpty())
			Expect(token.ExpiresAt).ToNot(BeZero())
		})
	})

	Describe("Validate", func() {
		It("accepts a no expired token and return user_id", func() {
			var (
				expectedUserID         = uuid.New()
				expectedExpirationDate = time.Duration(1 * time.Hour)
				tokenDetails           = auth.AccessToken{
					Expiration: expectedExpirationDate,
					UserID:     expectedUserID,
				}
				holder auth.AccessToken
			)
			token, _ := api.CreateAccessToken(&tokenDetails)

			log.Println(holder)

			Expect(api.ValidateToken(context.Background(), token.Token, &holder)).To(Succeed())
			Expect(holder.UserID).To(Equal(expectedUserID))
		})

		When("token is expired", func() {
			It("returns an error", func() {
				var (
					token  = "v4.local.hzqmiwdzwXasj7HJF5JbIQ3bvsg1Ph6cSvYXBTq-K_Uu3gHoJkcn7sUJJzJzuFLwcSk4Bg8taEPjHLulIQteRh9bq3ltrVir4M7B_O0fVg4OqNpO-htyMqse6CcjClIPpAQDDn_qH5_l3vb-ovqk7LXqapOCu20zIMwtNuwkRyv7xWR-8wQcN1_z64eSDe0mKloHLLi6ovRfZ1lpsGHmDZOtgzjHXPYkkLKpl3Fv2-BQBhPI8Lo"
					holder auth.AccessToken
				)

				Expect(api.ValidateToken(context.Background(), token, &holder)).ToNot(Succeed())
			})
		})
	})

})
