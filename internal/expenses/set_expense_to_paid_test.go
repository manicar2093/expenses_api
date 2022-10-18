package expenses_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetExpenseToPaid", func() {
	var (
		expensesRepoMock *mocks.ExpensesRepository
		timeGetableMock  *mocks.TimeGetable
		ctx              context.Context
		api              *expenses.ExpenseServiceImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		timeGetableMock = &mocks.TimeGetable{}
		ctx = context.Background()
		api = expenses.NewExpenseServiceImpl(expensesRepoMock, timeGetableMock)
	})

	AfterEach(func() {
		expensesRepoMock.AssertExpectations(GinkgoT())
	})

	It("change expense to paid", func() {
		var (
			expectedID = faker.Name()
			req        = expenses.SetExpenseToPaidInput{
				ID: expectedID,
			}
			expectedStatusCallUpdateExpenseStatus = true
		)
		expensesRepoMock.EXPECT().UpdateIsPaidByExpenseID(
			ctx,
			expectedID,
			expectedStatusCallUpdateExpenseStatus,
		).Return(nil)

		err := api.SetToPaid(ctx, &req)

		Expect(err).ToNot(HaveOccurred())
	})
})
