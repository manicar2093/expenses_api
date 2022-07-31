package recurrentexpenses_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/mocks"
)

var _ = Describe("GetAll", func() {
	var (
		recurrentExpensesRepoMock *mocks.RecurrentExpenseRepo
		ctx                       context.Context
		api                       *recurrentexpenses.GetAllRecurrentExpensesImpl
	)
	BeforeEach(func() {
		recurrentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		ctx = context.Background()
		api = recurrentexpenses.NewGetAllRecurrentExpensesImpl(recurrentExpensesRepoMock)
	})

	AfterEach(func() {
		recurrentExpensesRepoMock.AssertExpectations(GinkgoT())
	})

	It("returns all registered recurrent expenses", func() {
		var (
			expectedRepoReturn = []entities.RecurrentExpense{
				{}, {}, {},
			}
			expectedRecurrentExpensesCount = uint(len(expectedRepoReturn))
		)
		recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(&expectedRepoReturn, nil)

		got, err := api.GetAll(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.RecurrenteExpensesCount).To(Equal(expectedRecurrentExpensesCount))
		Expect(got.RecurrentExpenses).To(Equal(expectedRepoReturn))
	})
})
