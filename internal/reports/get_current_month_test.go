package reports_test

import (
	"context"
	"time"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetCurrentMonth", func() {

	var (
		expensesRepoMock *mocks.ExpensesRepository
		timeGetterMock   *mocks.TimeGetable
		timeGetterReturn time.Time
		ctx              context.Context
		service          *reports.CurrentMonthDetails
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		timeGetterMock = &mocks.TimeGetable{}
		timeGetterReturn = time.Date(2022, time.July, 1, 0, 0, 0, 0, time.Local)
		timeGetterMock.EXPECT().GetCurrentTime().Return(timeGetterReturn)
		ctx = context.Background()
		service = reports.NewCurrentMonthDetailsImpl(expensesRepoMock, timeGetterMock)
	})

	It("returns expenses data", func() {
		// should return expenses quantity and total amount
		var (
			expectedAmount1    = 231.90
			expectedAmount2    = 123.90
			expectedAmount3    = 321.90
			expectedRepoReturn = []entities.Expense{
				{Amount: expectedAmount1, Month: uint(time.July)},
				{Amount: expectedAmount2, Month: uint(time.July)},
				{Amount: expectedAmount3, Month: uint(time.July)},
			}
			expectedExpensesCount = uint(len(expectedRepoReturn))
			expectedTotalAmount   = expectedAmount1 + expectedAmount2 + expectedAmount3
		)
		expensesRepoMock.EXPECT().GetExpensesByMonth(ctx, time.July).Return(&expectedRepoReturn, nil)

		got, err := service.GetExpenses(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.TotalAmount).To(Equal(expectedTotalAmount))
		Expect(got.TotalExpenses).To(Equal(expectedExpensesCount))
		Expect(got.Expenses).To(Equal(&expectedRepoReturn))
	})

})
