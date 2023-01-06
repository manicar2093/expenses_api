package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/controllers"
	"github.com/manicar2093/expenses_api/internal/auth"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)

var _ = Describe("/auth", func() {

	var (
		e                          *echo.Echo
		googleLoginableByTokenMock *mocks.LoginableByToken
		api                        *controllers.LoginController
	)

	BeforeEach(func() {
		T := GinkgoT()
		e = echo.New()
		googleLoginableByTokenMock = mocks.NewLoginableByToken(T)
		api = controllers.NewLoginController(googleLoginableByTokenMock, e)
	})

	Describe("/login/google", func() {
		It("login user with given user google token", func() {
			var (
				expectedToken         = faker.Paragraph()
				expectedUserAgent     = faker.Name()
				expectedRemoteAddress = faker.IPv4()
				expectedJsonData      = fmt.Sprintf(`
				{"token": "%v"}`,
					expectedToken,
				)
				expectedLoginInput = auth.LoginInput{
					Token:     expectedToken,
					UserAgent: expectedUserAgent,
					ClientIP:  expectedRemoteAddress,
				}
				expectedLoginOutput = auth.LoginOutput{
					AccessToken:          faker.Name(),
					AccessTokenExpiresAt: time.Now(),
					RefreshToken:         uuid.New(),
				}
				req = testfunc.CreateJsonRequestForTest(http.MethodPost, "/expenses/", strings.NewReader(expectedJsonData))
				rec = httptest.NewRecorder()

				ctx = e.NewContext(req, rec)
			)
			ctx.Request().Header.Set("User-Agent", expectedUserAgent)
			ctx.Request().RemoteAddr = expectedRemoteAddress

			googleLoginableByTokenMock.EXPECT().Login(ctx.Request().Context(), &expectedLoginInput).Return(
				&expectedLoginOutput,
				nil,
			)

			err := api.LoginWGoogle(ctx)
			var body map[string]interface{}

			Expect(json.Unmarshal(rec.Body.Bytes(), &body)).To(Succeed())
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(body).To(gstruct.MatchAllKeys(gstruct.Keys{
				"access_token":            Equal(expectedLoginOutput.AccessToken),
				"access_token_expires_at": Not(BeZero()),
				"refresh_token":           Equal(expectedLoginOutput.RefreshToken.String()),
			}))
		})
	})

})
