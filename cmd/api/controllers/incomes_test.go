package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/cmd/api/controllers"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/incomes"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("/incomes", func() {

	var (
		e                  *echo.Echo
		incomesCreatorMock *mocks.CreateIncome
		middlewaresMock    *mocks.Middlewares
		api                *controllers.IncomesController
	)

	BeforeEach(func() {
		T := GinkgoT()
		e = testfunc.EchoWithValidator()
		incomesCreatorMock = mocks.NewCreateIncome(T)
		middlewaresMock = mocks.NewMiddlewares(T)
		api = controllers.NewIncomesController(middlewaresMock, incomesCreatorMock, e)
	})

	Describe("/", func() {
		When("POST", func() {
			It("create a new income", func() {
				var (
					expectedUserID                = uuid.New()
					expectedIncomeCreateInputCall = incomes.CreateIncomeInput{
						Income: entities.Income{
							UserID:      goption.Of(expectedUserID),
							Name:        faker.Name(),
							Amount:      12.8,
							Description: faker.Paragraph(),
						},
					}
					incomeJsonData = fmt.Sprintf(`{
						"name": "%v",
						"amount": %v,
						"description": "%v"
					}`, expectedIncomeCreateInputCall.Name, expectedIncomeCreateInputCall.Amount, expectedIncomeCreateInputCall.Description)
					req = testfunc.CreateJsonRequestForTest(http.MethodPost, "/incomes/", strings.NewReader(incomeJsonData))
					rec = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				ctx.Set("user_id", expectedUserID.String())
				incomesCreatorMock.EXPECT().Create(ctx.Request().Context(), &expectedIncomeCreateInputCall).Return(&expectedIncomeCreateInputCall.Income, nil)

				err := api.Create(ctx)
				var body map[string]interface{}

				Expect(json.Unmarshal(rec.Body.Bytes(), &body)).To(Succeed())
				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusCreated))
			})
		})
	})

})
