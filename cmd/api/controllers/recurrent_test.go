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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"gopkg.in/guregu/null.v4"

	"github.com/manicar2093/expenses_api/cmd/api/controllers"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("/recurrent_expenses", func() {
	var (
		e                              *echo.Echo
		createRecurrentExpense         *mocks.RecurrentExpenseCreatable
		getAllRecurrentExpenses        *mocks.RecurrentExpensesAllGettable
		createMonthlyRecurrentExpenses *mocks.MonthlyRecurrentExpensesCreateable
		middlewaresMock                *mocks.Middlewares
		api                            *controllers.RecurrentExpensesController
	)

	BeforeEach(func() {
		e = testfunc.EchoWithValidator()
		createRecurrentExpense = &mocks.RecurrentExpenseCreatable{}
		getAllRecurrentExpenses = &mocks.RecurrentExpensesAllGettable{}
		createMonthlyRecurrentExpenses = &mocks.MonthlyRecurrentExpensesCreateable{}
		middlewaresMock = &mocks.Middlewares{}
		api = controllers.NewRecurrentExpensesController(
			createRecurrentExpense,
			getAllRecurrentExpenses,
			createMonthlyRecurrentExpenses,
			middlewaresMock,
			e,
		)
	})

	AfterEach(func() {
		T := GinkgoT()
		createRecurrentExpense.AssertExpectations(T)
		getAllRecurrentExpenses.AssertExpectations(T)
		createMonthlyRecurrentExpenses.AssertExpectations(T)
		middlewaresMock.AssertExpectations(T)
	})

	Describe("/", func() {
		When("POST", func() {
			It("creates a new recurrent expense", func() {
				var (
					expectedUserIDUUID                  = uuid.New()
					expectedUserID                      = expectedUserIDUUID.String()
					expectedCreateRecurrentExpenseInput = recurrentexpenses.CreateRecurrentExpenseInput{
						Name:        faker.Name(),
						Amount:      12.8,
						Description: faker.Paragraph(),
						UserID:      expectedUserID,
					}
					expectedCreateExpenseOutput = recurrentexpenses.CreateRecurrentExpenseOutput{
						RecurrentExpense: &entities.RecurrentExpense{
							ID:          uuid.New(),
							UserID:      expectedUserIDUUID,
							Name:        expectedCreateRecurrentExpenseInput.Name,
							Amount:      expectedCreateRecurrentExpenseInput.Amount,
							Description: null.StringFrom(expectedCreateRecurrentExpenseInput.Description),
							CreatedAt:   testfunc.ToPointer(time.Now()),
							UpdatedAt:   testfunc.ToPointer(time.Now()),
						},
						NextMonthExpense: &entities.Expense{
							ID:          uuid.New(),
							UserID:      expectedUserIDUUID,
							Name:        null.StringFrom(expectedCreateRecurrentExpenseInput.Name),
							Amount:      expectedCreateRecurrentExpenseInput.Amount,
							Description: null.StringFrom(expectedCreateRecurrentExpenseInput.Description),
							Day:         1,
							Month:       1,
							Year:        1,
							IsPaid:      false,
							CreatedAt:   testfunc.ToPointer(time.Now()),
							UpdatedAt:   testfunc.ToPointer(time.Now()),
						},
					}
					expectedJson = fmt.Sprintf(`{"name": "%v","amount": %v,"description": "%v"}`,
						expectedCreateRecurrentExpenseInput.Name,
						expectedCreateRecurrentExpenseInput.Amount,
						expectedCreateRecurrentExpenseInput.Description,
					)
					req = testfunc.CreateJsonRequestForTest(http.MethodPost, "/recurrent_expenses/", strings.NewReader(expectedJson))
					rec = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				ctx.Set("user_id", expectedUserID)
				createRecurrentExpense.EXPECT().CreateRecurrentExpense(
					ctx.Request().Context(),
					&expectedCreateRecurrentExpenseInput,
				).Return(
					&expectedCreateExpenseOutput,
					nil,
				)

				err := api.Create(ctx)
				var body map[string]interface{}

				Expect(err).ToNot(HaveOccurred())
				Expect(json.Unmarshal(rec.Body.Bytes(), &body)).To(Succeed())
				Expect(rec.Code).To(Equal(http.StatusCreated))
				Expect(body).To(gstruct.MatchAllKeys(gstruct.Keys{
					"recurrent_expense": gstruct.MatchAllKeys(
						gstruct.Keys{
							"id":          Not(BeEmpty()),
							"user_id":     Equal(expectedUserID),
							"name":        Equal(expectedCreateRecurrentExpenseInput.Name),
							"amount":      Equal(expectedCreateRecurrentExpenseInput.Amount),
							"description": Equal(expectedCreateRecurrentExpenseInput.Description),
							"created_at":  Not(BeEmpty()),
							"updated_at":  Not(BeEmpty()),
						},
					),
					"next_month_expense": gstruct.MatchAllKeys(
						gstruct.Keys{
							"id":                   Not(BeEmpty()),
							"recurrent_expense_id": BeNil(),
							"user_id":              Equal(expectedUserID),
							"name":                 Equal(expectedCreateRecurrentExpenseInput.Name),
							"amount":               Equal(expectedCreateRecurrentExpenseInput.Amount),
							"description":          Equal(expectedCreateRecurrentExpenseInput.Description),
							"day":                  Not(BeZero()),
							"month":                Not(BeZero()),
							"year":                 Not(BeZero()),
							"is_paid":              BeFalse(),
							"created_at":           Not(BeEmpty()),
							"updated_at":           Not(BeEmpty()),
						},
					),
				}))
			})
		})
	})

	Describe("/all", func() {
		When("GET", func() {
			It("returns all recurrent expenses", func() {
				var (
					expectedUserIDUUID = uuid.New()
					expectedUserID     = expectedUserIDUUID.String()
					req                = testfunc.CreateJsonRequestForTest(http.MethodGet, "/recurrent_expenses/all", nil)
					rec                = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				ctx.Set("user_id", expectedUserID)
				getAllRecurrentExpenses.EXPECT().GetAll(ctx.Request().Context(), expectedUserIDUUID).Return(
					&recurrentexpenses.GetAllRecurrentExpensesOutput{
						RecurrentExpenses:       []*entities.RecurrentExpense{{}},
						RecurrenteExpensesCount: 1,
					},
					nil,
				)

				err := api.GetAll(ctx)

				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("/monthly_expenses", func() {
		When("POST", func() {
			It("creates all recurrent expenses for the month", func() {
				var (
					expectedUserIDUUID = uuid.New()
					expectedUserID     = expectedUserIDUUID.String()
					req                = testfunc.CreateJsonRequestForTest(http.MethodGet, "/recurrent_expenses/monthly_expenses", nil)
					rec                = httptest.NewRecorder()

					ctx = e.NewContext(req, rec)
				)
				ctx.Set("user_id", expectedUserID)
				createMonthlyRecurrentExpenses.EXPECT().CreateMonthlyRecurrentExpenses(ctx.Request().Context(), expectedUserIDUUID).Return(
					&recurrentexpenses.CreateMonthlyRecurrentExpensesOutput{},
					nil,
				)

				err := api.CreateMonthly(ctx)

				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
