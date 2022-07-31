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

	AfterEach(func() {
		T := GinkgoT()
		expensesRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
	})

	It("returns expenses data", func() {
		// should return expenses quantity and total amount
		var (
			expectedPaidAmount1  = 231.90
			expectedPaidAmount2  = 123.90
			expectedPaidAmount3  = 321.90
			expectedPaidExpenses = []entities.Expense{
				{Amount: expectedPaidAmount1, Month: uint(time.July), IsRecurrent: false, IsPaid: true},
				{Amount: expectedPaidAmount2, Month: uint(time.July), IsRecurrent: false, IsPaid: true},
				{Amount: expectedPaidAmount3, Month: uint(time.July), IsRecurrent: true, IsPaid: true},
			}
			expectedUnpaidAmount1  = 234.90
			expectedUnpaidAmount2  = 345.90
			expectedUnpaidExpenses = []entities.Expense{
				{Amount: expectedUnpaidAmount1, Month: uint(time.July), IsRecurrent: true, IsPaid: false},
				{Amount: expectedUnpaidAmount2, Month: uint(time.July), IsRecurrent: true, IsPaid: false},
			}
			expectedTotalPaidAmount     = expectedPaidAmount1 + expectedPaidAmount2 + expectedPaidAmount3
			expectedTotalUnpaidAmount   = expectedUnpaidAmount1 + expectedUnpaidAmount2
			expectedRepoReturn          = append(expectedPaidExpenses, expectedUnpaidExpenses...)
			expectedTotalExpenses       = uint(len(expectedRepoReturn))
			expectedPaidExpensesCount   = uint(len(expectedPaidExpenses))
			expectedUnpaidExpensesCount = uint(len(expectedUnpaidExpenses))
		)
		expensesRepoMock.EXPECT().GetExpensesByMonth(ctx, time.July).Return(&expectedRepoReturn, nil)

		got, err := service.GetExpenses(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.TotalPaidAmount).To(Equal(expectedTotalPaidAmount))
		Expect(got.TotalUnpaidAmount).To(Equal(expectedTotalUnpaidAmount))
		Expect(got.ExpensesCount).To(Equal(expectedTotalExpenses))
		Expect(got.PaidExpensesCount).To(Equal(expectedPaidExpensesCount))
		Expect(got.UnpaidExpensesCount).To(Equal(expectedUnpaidExpensesCount))
		Expect(got.Expenses).To(Equal(expectedRepoReturn))
		Expect(got.PaidExpenses).To(Equal(expectedPaidExpenses))
		Expect(got.UnpaidExpenses).To(Equal(expectedUnpaidExpenses))
	})

})
