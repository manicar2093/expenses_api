package expenses_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/internal/schemas"
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
		expensesRepoMock = &mocks.ExpensesRepository{}
		timeGetableMock = &mocks.TimeGetable{}
		ctx = context.Background()
		service = expenses.NewExpenseServiceImpl(expensesRepoMock, timeGetableMock)
	})

	It("toggles expense IsPaid status", func() {
		var (
			expectedExpenseID         = primitive.NewObjectID()
			expectedExpenseIDString   = expectedExpenseID.Hex()
			expectedExpenseWithStatus = schemas.ExpenseIDWithIsPaidStatus{
				ID:     expectedExpenseID,
				IsPaid: true,
			}
			expectedIsPaidUpdateCall   = !expectedExpenseWithStatus.IsPaid
			expectedToggleExpenseInput = expenses.ToggleExpenseIsPaidInput{
				ID: expectedExpenseIDString,
			}
		)
		expensesRepoMock.EXPECT().GetExpenseStatusByID(
			ctx,
			expectedExpenseIDString,
		).Return(&expectedExpenseWithStatus, nil)
		expensesRepoMock.EXPECT().UpdateIsPaidByExpenseID(
			ctx,
			expectedExpenseIDString,
			expectedIsPaidUpdateCall,
		).Return(nil)

		got, err := service.ToggleIsPaid(ctx, &expectedToggleExpenseInput)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.ID).To(Equal(expectedExpenseIDString))
		Expect(got.CurrentIsPaidStatus).To(Equal(expectedIsPaidUpdateCall))

	})

})
