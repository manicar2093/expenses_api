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
	"github.com/onsi/gomega/gstruct"

	"github.com/manicar2093/expenses_api/cmd/api/controllers"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("/expenses", func() {

	var (
		e                      *echo.Echo
		createableExpensesMock *mocks.ExpenseCreatable
		toPaidSetteable        *mocks.ExpenseToPaidSetteable
		isPaidToggleable       *mocks.ExpenseToPaidTogglable
		expenseUpdateable      *mocks.ExpenseUpdateable
		middlewaresMock        *mocks.Middlewares
		api                    controllers.ExpensesController
	)

	BeforeEach(func() {
		T := GinkgoT()
		e = testfunc.EchoWithValidator()
		createableExpensesMock = mocks.NewExpenseCreatable(T)
		toPaidSetteable = mocks.NewExpenseToPaidSetteable(T)
		isPaidToggleable = mocks.NewExpenseToPaidTogglable(T)
		expenseUpdateable = mocks.NewExpenseUpdateable(T)
		middlewaresMock = mocks.NewMiddlewares(T)
		api = *controllers.NewExpensesController(
			createableExpensesMock,
			toPaidSetteable,
			isPaidToggleable,
			expenseUpdateable,
			middlewaresMock,
			e,
		)
	})

	Describe("/", func() {
		When("POST", func() {
			It("creates a new expense", func() {
				var (
					expectedExpenseCreated = testfunc.GeneratePaidExpense()
					expectedUserID         = expectedExpenseCreated.UserID.String()
					expectedExpenseCall    = expenses.CreateExpenseInput{
						Name:        expectedExpenseCreated.Name.String,
						Amount:      expectedExpenseCreated.Amount,
						Description: expectedExpenseCreated.Description.String,
						UserID:      expectedUserID,
					}
					expensesData = fmt.Sprintf(`
					{
						"name": "%v",
						"amount": %v,
						"description": "%v"
					}`, expectedExpenseCall.Name, expectedExpenseCall.Amount, expectedExpenseCall.Description)
					req = testfunc.CreateJsonRequestForTest(http.MethodPost, "/expenses/", strings.NewReader(expensesData))
					rec = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				ctx.Set("user_id", expectedUserID)
				createableExpensesMock.EXPECT().CreateExpense(ctx.Request().Context(), &expectedExpenseCall).Return(expectedExpenseCreated, nil)

				err := api.Create(ctx)
				var body map[string]interface{}

				Expect(json.Unmarshal(rec.Body.Bytes(), &body)).To(Succeed())
				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusCreated))
				Expect(body).To(gstruct.MatchAllKeys(gstruct.Keys{
					"id":                   Not(BeEmpty()),
					"recurrent_expense_id": BeNil(),
					"user_id":              Equal(expectedUserID),
					"name":                 Equal(expectedExpenseCreated.Name.String),
					"amount":               Equal(expectedExpenseCreated.Amount),
					"description":          Equal(expectedExpenseCreated.Description.String),
					"is_paid":              BeTrue(),
					"created_at":           Not(BeEmpty()),
					"updated_at":           Not(BeEmpty()),
				}))
			})

		})
	})

	Describe("/to_paid", func() {
		When("PUT", func() {
			It("set a expense to paid", func() {
				var (
					expenseID             = uuid.New().String()
					expectedExpenseIDData = fmt.Sprintf(`
					{"id": "%v"}
					`, expenseID)
					req = testfunc.CreateJsonRequestForTest(http.MethodPut, "/expenses/to_paid", strings.NewReader(expectedExpenseIDData))
					rec = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				toPaidSetteable.EXPECT().SetToPaid(ctx.Request().Context(), &expenses.SetExpenseToPaidInput{ID: expenseID}).Return(nil)

				err := api.ToPaid(ctx)

				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.Bytes()).To(HaveLen(0))
			})
		})
	})

	Describe("/toogle_is_paid", func() {
		When("PUT", func() {
			It("toogles expense is paid by id", func() {
				var (
					expenseIDUUID         = uuid.New()
					expenseID             = expenseIDUUID.String()
					expectedExpenseIDData = fmt.Sprintf(`
					{"id": "%v"}
					`, expenseID)
					req = testfunc.CreateJsonRequestForTest(http.MethodPut, "/expenses/toogle_is_paid", strings.NewReader(expectedExpenseIDData))
					rec = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				isPaidToggleable.EXPECT().ToggleIsPaid(ctx.Request().Context(), &expenses.ToggleExpenseIsPaidInput{ID: expenseID}).Return(&expenses.ToggleExpenseIsPaidOutput{ID: expenseIDUUID, CurrentIsPaidStatus: true}, nil)

				err := api.ToggleIsPaid(ctx)
				var body map[string]interface{}

				Expect(json.Unmarshal(rec.Body.Bytes(), &body)).To(Succeed())
				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(body).To(gstruct.MatchAllKeys(gstruct.Keys{
					"id":                     Not(BeEmpty()),
					"current_is_paid_status": BeTrue(),
				}))
			})
		})
	})

	Describe("/update", func() {
		When("PUT", func() {
			It("updates a given expense", func() {
				var (
					expenseUpdateInput = expenses.UpdateExpenseInput{
						ID:          uuid.New().String(),
						Name:        faker.Name(),
						Amount:      12.3,
						Description: faker.Paragraph(),
					}
					expenseUpdateJson = fmt.Sprintf(`
					{
						"id": "%v",
						"name":"%v",
						"amount":%v,
						"description":"%v"
					}`, expenseUpdateInput.ID, expenseUpdateInput.Name, expenseUpdateInput.Amount, expenseUpdateInput.Description)
					req = testfunc.CreateJsonRequestForTest(http.MethodPut, "/expenses/update", strings.NewReader(expenseUpdateJson))
					rec = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				expenseUpdateable.EXPECT().UpdateExpense(ctx.Request().Context(), &expenseUpdateInput).Return(nil)

				Expect(api.Update(ctx)).To(Succeed())
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.Bytes()).To(HaveLen(0))
			})
		})
	})

})
