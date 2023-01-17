package expenses_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetExpenseToPaid", func() {
	var (
		expenseRepoMock *mocks.ExpensesRepository
		timeGetterMock  *mocks.TimeGetable
		ctx             context.Context
		api             *expenses.ExpenseServiceImpl
	)

	BeforeEach(func() {
		T := GinkgoT()
		expenseRepoMock = mocks.NewExpensesRepository(T)
		timeGetterMock = mocks.NewTimeGetable(T)
		ctx = context.Background()
		api = expenses.NewExpenseServiceImpl(expenseRepoMock, timeGetterMock)
	})

	It("change expense to paid", func() {
		var (
			expectedIDAsString = uuid.New().String()
			expectedID         = uuid.MustParse(expectedIDAsString)
			req                = expenses.SetExpenseToPaidInput{
				ID: expectedIDAsString,
			}
			expectedStatusCallUpdateExpenseStatus = true
		)
		expenseRepoMock.EXPECT().UpdateIsPaidByExpenseID(
			ctx,
			expectedID,
			expectedStatusCallUpdateExpenseStatus,
		).Return(nil)

		err := api.SetToPaid(ctx, &req)

		Expect(err).ToNot(HaveOccurred())
	})

})
