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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"

	"github.com/manicar2093/expenses_api/cmd/api/controllers"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/period"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("Recurrent", func() {

	var (
		e                           *echo.Echo
		recurrentExpenseCreatorMock *mocks.RecurrentExpenseCreator
		middlewaresMock             *mocks.Middlewares
		api                         *controllers.RecurrentExpensesController
	)

	BeforeEach(func() {
		T := GinkgoT()
		e = testfunc.EchoWithValidator()
		recurrentExpenseCreatorMock = mocks.NewRecurrentExpenseCreator(T)
		middlewaresMock = mocks.NewMiddlewares(T)
		api = controllers.NewRecurrentExpensesController(recurrentExpenseCreatorMock, middlewaresMock, e)
	})

	Describe("/", func() {
		When("POST", func() {
			It("creates a new recurrent expense", func() {
				var (
					expectedUserID           = uuid.New()
					expectedRecurrentExpense = entities.RecurrentExpense{
						Name:        faker.Name(),
						Amount:      1050.0,
						Description: null.StringFrom(faker.Paragraph()),
						Periodicity: period.BiMonthly,
						UserID:      expectedUserID,
					}
					jsonData = fmt.Sprintf(`
					{
						"name": "%v",
						"amount": %v,
						"description": "%v",
						"periodicity": "%v"
					}`, expectedRecurrentExpense.Name,
						expectedRecurrentExpense.Amount,
						expectedRecurrentExpense.Description.String,
						expectedRecurrentExpense.Periodicity,
					)
					req = testfunc.CreateJsonRequestForTest(http.MethodPost, "/recurrent_expense", strings.NewReader(jsonData))
					rec = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				ctx.Set("user_id", expectedUserID.String())
				recurrentExpenseCreatorMock.EXPECT().Create(ctx.Request().Context(), &expectedRecurrentExpense).Return(nil)

				Expect(api.Create(ctx)).To(Succeed())
				var body map[string]interface{}
				Expect(json.Unmarshal(rec.Body.Bytes(), &body)).To(Succeed())
				Expect(body).To(HaveKeyWithValue("name", expectedRecurrentExpense.Name))

			})
		})
	})
})
