package expenses_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/converters"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			expectedID = primitive.NewObjectID().Hex()
			req        = expenses.SetExpenseToPaidInput{
				ID: expectedID,
			}
			expectedIDCallUpdateExpenseStatus, _  = primitive.ObjectIDFromHex(expectedID)
			expectedStatusCallUpdateExpenseStatus = true
		)
		expensesRepoMock.EXPECT().UpdateIsPaidByExpenseID(
			ctx,
			expectedIDCallUpdateExpenseStatus,
			expectedStatusCallUpdateExpenseStatus,
		).Return(nil)

		err := api.SetToPaid(ctx, &req)

		Expect(err).ToNot(HaveOccurred())
	})

	When("given id it is not valid", func() {
		It("returns a handleable error", func() {
			var (
				expectedID = faker.Name()
				req        = expenses.SetExpenseToPaidInput{
					ID: expectedID,
				}
			)

			err := api.SetToPaid(ctx, &req)

			Expect(err).To(BeAssignableToTypeOf(&converters.IDNotValidIDError{}))
		})
	})
})
