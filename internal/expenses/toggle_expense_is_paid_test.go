package expenses_test

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

var _ = Describe("ToggleExpenseIsPaid", func() {

	var (
		expensesRepoMock *mocks.ExpensesRepository
		timeGetableMock  *mocks.TimeGetable
		ctx              context.Context
		validatorMock    *mocks.StructValidable
		service          *expenses.ExpenseServiceImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		timeGetableMock = &mocks.TimeGetable{}
		ctx = context.Background()
		validatorMock = &mocks.StructValidable{}
		service = expenses.NewExpenseServiceImpl(expensesRepoMock, timeGetableMock, validatorMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		expensesRepoMock.AssertExpectations(T)
		timeGetableMock.AssertExpectations(T)
		validatorMock.AssertExpectations(T)
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
		validatorMock.EXPECT().ValidateStruct(&expectedToggleExpenseInput).Return(nil)
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

	When("request is not valid", Label(testfunc.IntegrationTest), func() {
		It("return a validation error", func() {
			var invalidRequest = expenses.ToggleExpenseIsPaidInput{
				"not uuid",
			}
			integrationTestApi := expenses.NewExpenseServiceImpl(expensesRepoMock, timeGetableMock, validator.NewGooKitValidator())

			got, err := integrationTestApi.ToggleIsPaid(ctx, &invalidRequest)

			Expect(got).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

})
