package recurrentexpenses_test

import (
	"context"

	"github.com/google/uuid"
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
		userID                    uuid.UUID
		api                       *recurrentexpenses.RecurrentExpenseServiceImpl
	)
	BeforeEach(func() {
		recurrentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		ctx = context.Background()
		userID = uuid.New()
		api = recurrentexpenses.NewGetAllRecurrentExpenseServiceImpl(recurrentExpensesRepoMock)
	})

	AfterEach(func() {
		recurrentExpensesRepoMock.AssertExpectations(GinkgoT())
	})

	It("returns all registered recurrent expenses", func() {
		var (
			expectedRepoReturn = []*entities.RecurrentExpense{
				{}, {}, {},
			}
			expectedRecurrentExpensesCount = uint(len(expectedRepoReturn))
		)
		recurrentExpensesRepoMock.EXPECT().FindAll(ctx, userID).Return(expectedRepoReturn, nil)

		got, err := api.GetAll(ctx, userID)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.RecurrenteExpensesCount).To(Equal(expectedRecurrentExpensesCount))
		Expect(got.RecurrentExpenses).To(Equal(expectedRepoReturn))
	})
})
