package controllers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/cmd/api/controllers"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("/reports", func() {

	var (
		e               *echo.Echo
		getCurrentMonth *mocks.CurrentMonthDetailsGettable
		middlewaresMock *mocks.Middlewares
		api             *controllers.ReportsController
	)

	BeforeEach(func() {
		e = echo.New()
		getCurrentMonth = &mocks.CurrentMonthDetailsGettable{}
		middlewaresMock = &mocks.Middlewares{}
		api = controllers.NewReportsController(getCurrentMonth, middlewaresMock, e)

	})

	AfterEach(func() {
		T := GinkgoT()
		getCurrentMonth.AssertExpectations(T)
		middlewaresMock.AssertExpectations(T)
	})

	Describe("/current_month", func() {
		When("GET", func() {
			It("retreives month report", func() {
				var (
					expectedUserIDUUID = uuid.New()
					expectedUserID     = expectedUserIDUUID.String()
					req                = testfunc.CreateJsonRequestForTest(http.MethodPost, "/expenses/", nil)
					rec                = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				ctx.Set("user_id", expectedUserID)
				getCurrentMonth.EXPECT().GetExpenses(ctx.Request().Context(), expectedUserIDUUID).Return(
					&reports.CurrentMonthDetailsOutput{},
					nil,
				)

				err := api.CurrentMonth(ctx)

				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.Bytes()).ToNot(BeEmpty())
			})
		})
	})

})
