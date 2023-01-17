package expenses_test

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
)

var _ = Describe("ToggleExpenseIsPaid", func() {

	var (
		expensesRepoMock *mocks.ExpensesRepository
		timeGetableMock  *mocks.TimeGetable
		ctx              context.Context
		service          *expenses.ExpenseServiceImpl
	)

	BeforeEach(func() {
		T := GinkgoT()
		expensesRepoMock = mocks.NewExpensesRepository(T)
		timeGetableMock = mocks.NewTimeGetable(T)
		ctx = context.Background()
		service = expenses.NewExpenseServiceImpl(expensesRepoMock, timeGetableMock)
	})

	It("toggles expense IsPaid status", func() {
		var (
			expectedExpenseIDAsString = uuid.New().String()
			expectedExpenseID         = uuid.MustParse(expectedExpenseIDAsString)
			expectedExpenseWithStatus = entities.ExpenseIDWithIsPaidStatus{
				ID:     expectedExpenseID,
				IsPaid: true,
			}
			expectedIsPaidUpdateCall   = !expectedExpenseWithStatus.IsPaid
			expectedToggleExpenseInput = expenses.ToggleExpenseIsPaidInput{
				ID: expectedExpenseIDAsString,
			}
		)
		expensesRepoMock.EXPECT().GetExpenseStatusByID(
			ctx,
			expectedExpenseID,
		).Return(&expectedExpenseWithStatus, nil)
		expensesRepoMock.EXPECT().UpdateIsPaidByExpenseID(
			ctx,
			expectedExpenseID,
			expectedIsPaidUpdateCall,
		).Return(nil)

		got, err := service.ToggleIsPaid(ctx, &expectedToggleExpenseInput)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.ID).To(Equal(expectedExpenseID))
		Expect(got.CurrentIsPaidStatus).To(Equal(expectedIsPaidUpdateCall))

	})

})
