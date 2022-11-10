package recurrentexpenses_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {

	var (
		recurentExpensesRepoMock *mocks.RecurrentExpenseRepo
		ctx                      context.Context
		api                      *recurrentexpenses.CreateRecurrentExpenseImpl
	)

	BeforeEach(func() {
		recurentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		ctx = context.Background()
		api = recurrentexpenses.NewCreateRecurrentExpenseImpl(recurentExpensesRepoMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		recurentExpensesRepoMock.AssertExpectations(T)
	})

	It("creates recurrent expense with next month expense record", func() {
		var (
			expectedExpenseName        = faker.Name()
			expectedExpenseAmount      = faker.Latitude()
			expectedExpenseDescription = faker.Paragraph()
			// expectedExpensePeriodicity = periodtypes.Paydaily
			// expectedMonthToBeConsideredSince = 12
			request = recurrentexpenses.CreateRecurrentExpenseInput{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: expectedExpenseDescription,
				// Periodicity:                  expectedExpensePeriodicity,
				// FromMonthToBeConsideredSince: expectedMonthToBeConsideredSince,
			}
			expectedRecurrentExpenseSaved = entities.RecurrentExpense{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: expectedExpenseDescription,
				// Periodicity: expectedExpensePeriodicity,
			}
		)
		recurentExpensesRepoMock.EXPECT().Save(ctx, &expectedRecurrentExpenseSaved).Return(nil)

		got, err := api.Create(ctx, &request)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.RecurrentExpense).To(Equal(&expectedRecurrentExpenseSaved))
	})

})
