package expenses_test

import (
	"context"
	"fmt"
	"time"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
)

var _ = Describe("CreateImpl", func() {

	var (
		expenseRepoMock *mocks.ExpensesRepository
		timeGetterMock  *mocks.TimeGetable
		ctx             context.Context
		api             *expenses.ExpenseServiceImpl
	)

	BeforeEach(func() {
		expenseRepoMock = &mocks.ExpensesRepository{}
		timeGetterMock = &mocks.TimeGetable{}
		ctx = context.TODO()
		api = expenses.NewExpenseServiceImpl(expenseRepoMock, timeGetterMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		expenseRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
	})

	It("creates a new expense from schema", func() {
		var (
			expectedName              = faker.Name()
			expectedDescription       = faker.Paragraph()
			expectedAmount            = faker.Latitude()
			expectedCurrentDateReturn = time.Date(2022, time.August, 1, 0, 0, 0, 0, time.Local)
			request                   = expenses.CreateExpenseInput{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
			}
			expectedExpenseToSave = entities.Expense{
				Name:        null.StringFrom(expectedName),
				Amount:      expectedAmount,
				Day:         uint(expectedCurrentDateReturn.Day()),
				Month:       uint(expectedCurrentDateReturn.Month()),
				Year:        uint(expectedCurrentDateReturn.Year()),
				Description: null.StringFrom(expectedDescription),
				IsPaid:      true,
				CreatedAt:   &expectedCurrentDateReturn,
			}
		)
		timeGetterMock.EXPECT().GetCurrentTime().Return(expectedCurrentDateReturn)
		expenseRepoMock.EXPECT().Save(ctx, &expectedExpenseToSave).Return(nil)

		got, err := api.CreateExpense(ctx, &request)

		Expect(err).ToNot(HaveOccurred())
		Expect(got).To(Equal(&expectedExpenseToSave))
	})

	When("expense is asked to be created for next month", func() {
		It("assign the need date to be created", func() {
			var (
				expectedName                = faker.Name()
				expectedDescription         = faker.Paragraph()
				expectedAmount              = faker.Latitude()
				expectedNextMonthDateReturn = time.Date(2022, time.August, 1, 0, 0, 0, 0, time.Local)
				expectedNowDateReturn       = time.Date(2022, time.July, 30, 0, 0, 0, 0, time.Local)
				expectedDateString          = "Fecha de registro: 30/07/2022"
				expectedExpenseDescription  = fmt.Sprintf("%s\n\n%s", expectedDescription, expectedDateString)
				request                     = expenses.CreateExpenseInput{
					Name:         expectedName,
					Amount:       expectedAmount,
					Description:  expectedDescription,
					ForNextMonth: true,
				}
				expectedExpenseToSave = entities.Expense{
					Name:        null.StringFrom(expectedName),
					Amount:      expectedAmount,
					Day:         uint(expectedNextMonthDateReturn.Day()),
					Month:       uint(expectedNextMonthDateReturn.Month()),
					Year:        uint(expectedNextMonthDateReturn.Year()),
					Description: null.StringFrom(expectedExpenseDescription),
					IsPaid:      true,
					CreatedAt:   &expectedNextMonthDateReturn,
				}
			)
			expenseRepoMock.EXPECT().Save(ctx, &expectedExpenseToSave).Return(nil)
			timeGetterMock.EXPECT().GetCurrentTime().Return(expectedNowDateReturn)
			timeGetterMock.EXPECT().GetNextMonthAtFirtsDay().Return(expectedNextMonthDateReturn)

			got, err := api.CreateExpense(ctx, &request)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(Equal(&expectedExpenseToSave))
		})
	})

})
