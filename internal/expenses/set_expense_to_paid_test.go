package expenses_test

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/json"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	"github.com/manicar2093/expenses_api/pkg/validator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetExpenseToPaid", func() {
	var (
		expenseRepoMock *mocks.ExpensesRepository
		timeGetterMock  *mocks.TimeGetable
		ctx             context.Context
		validatorMock   *mocks.StructValidable
		api             *expenses.ExpenseServiceImpl
	)

	BeforeEach(func() {
		expenseRepoMock = &mocks.ExpensesRepository{}
		timeGetterMock = &mocks.TimeGetable{}
		ctx = context.Background()
		validatorMock = &mocks.StructValidable{}
		api = expenses.NewExpenseServiceImpl(expenseRepoMock, timeGetterMock, validatorMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		expenseRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
		validatorMock.AssertExpectations(T)
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
		validatorMock.EXPECT().ValidateStruct(&req).Return(nil)
		expenseRepoMock.EXPECT().UpdateIsPaidByExpenseID(
			ctx,
			expectedID,
			expectedStatusCallUpdateExpenseStatus,
		).Return(nil)

		err := api.SetToPaid(ctx, &req)

		Expect(err).ToNot(HaveOccurred())
	})

	When("request is not valid", Label(testfunc.IntegrationTest), func() {
		It("return a validation error", func() {
			var invalidRequest = expenses.SetExpenseToPaidInput{
				"not uuid",
			}

			integrationTestApi := expenses.NewExpenseServiceImpl(expenseRepoMock, timeGetterMock, validator.NewGooKitValidator())

			err := integrationTestApi.SetToPaid(ctx, &invalidRequest)

			log.Println(json.MustMarshall(err))
			Expect(err).To(HaveOccurred())
		})
	})
})
