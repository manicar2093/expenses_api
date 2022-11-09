package reports_test

import (
	"context"
	"time"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/reports"
	"github.com/manicar2093/expenses_api/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("GetCurrentMonth", func() {

	var (
		expensesRepoMock   *mocks.ExpensesRepository
		timeGetterMock     *mocks.TimeGetable
		timeGetterReturn   time.Time
		periodicityGenMock *mocks.ExpensePeriodicityCreateable
		ctx                context.Context
		service            *reports.CurrentMonthDetails
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		timeGetterMock = &mocks.TimeGetable{}
		timeGetterReturn = time.Date(2022, time.July, 1, 0, 0, 0, 0, time.Local)
		timeGetterMock.EXPECT().GetCurrentTime().Return(timeGetterReturn)
		periodicityGenMock = &mocks.ExpensePeriodicityCreateable{}
		ctx = context.Background()
		service = reports.NewCurrentMonthDetailsImpl(expensesRepoMock, timeGetterMock, periodicityGenMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		expensesRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
		periodicityGenMock.AssertExpectations(T)
	})

	It("returns expenses data", func() {
		var (
			expectedMonth        = uint(timeGetterReturn.Month())
			expectedYear         = uint(timeGetterReturn.Year())
			expectedPaidAmount1  = 231.90
			expectedPaidAmount2  = 123.90
			expectedPaidAmount3  = 321.90
			expectedPaidExpenses = []*entities.Expense{
				{Amount: expectedPaidAmount1, Month: expectedMonth, IsRecurrent: false, IsPaid: true},
				{Amount: expectedPaidAmount2, Month: expectedMonth, IsRecurrent: false, IsPaid: true},
				{Amount: expectedPaidAmount3, Month: expectedMonth, IsRecurrent: true, IsPaid: true},
			}
			expectedUnpaidAmount1  = 234.90
			expectedUnpaidAmount2  = 345.90
			expectedUnpaidExpenses = []*entities.Expense{
				{Amount: expectedUnpaidAmount1, Month: expectedMonth, IsRecurrent: true, IsPaid: false},
				{Amount: expectedUnpaidAmount2, Month: expectedMonth, IsRecurrent: true, IsPaid: false},
			}
			expectedTotalPaidAmount                 = expectedPaidAmount1 + expectedPaidAmount2 + expectedPaidAmount3
			expectedTotalUnpaidAmount               = expectedUnpaidAmount1 + expectedUnpaidAmount2
			expectedRepoReturn                      = append(expectedPaidExpenses, expectedUnpaidExpenses...)
			expectedTotalExpenses                   = uint(len(expectedRepoReturn))
			expectedPaidExpensesCount               = uint(len(expectedPaidExpenses))
			expectedUnpaidExpensesCount             = uint(len(expectedUnpaidExpenses))
			expectedRecurrentExpensesMonthlyCreated = entities.RecurrentExpensesMonthlyCreated{
				ID:    primitive.NewObjectID(),
				Month: expectedMonth,
				Year:  uint(expectedYear),
				ExpensesCount: []*entities.ExpensesCount{
					{RecurrentExpenseID: primitive.NewObjectID(), RecurrentExpense: &entities.RecurrentExpense{}, ExpensesRelatedIDs: []primitive.ObjectID{primitive.NewObjectID()}, TotalExpenses: 1, TotalExpensesPaid: 0},
				},
			}
		)
		expensesRepoMock.EXPECT().GetExpensesByMonth(ctx, time.July).Return(expectedRepoReturn, nil)
		periodicityGenMock.EXPECT().GenerateRecurrentExpensesByYearAndMonth(ctx, expectedMonth, expectedYear).Return(&expectedRecurrentExpensesMonthlyCreated, nil)

		got, err := service.GetCurrentMonthDetails(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.TotalPaidAmount).To(Equal(expectedTotalPaidAmount))
		Expect(got.TotalUnpaidAmount).To(Equal(expectedTotalUnpaidAmount))
		Expect(got.ExpensesCount).To(Equal(expectedTotalExpenses))
		Expect(got.PaidExpensesCount).To(Equal(expectedPaidExpensesCount))
		Expect(got.UnpaidExpensesCount).To(Equal(expectedUnpaidExpensesCount))
		Expect(got.Expenses).To(Equal(expectedRepoReturn))
		Expect(got.PaidExpenses).To(Equal(expectedPaidExpenses))
		Expect(got.UnpaidExpenses).To(Equal(expectedUnpaidExpenses))
		Expect(got.RecurrentExpensesDetails).To(Equal(&expectedRecurrentExpensesMonthlyCreated))
	})

})
