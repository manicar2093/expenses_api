package recurrentexpenses_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {

	var (
		expensesRepoMock         *mocks.ExpensesRepository
		recurentExpensesRepoMock *mocks.RecurrentExpenseRepo
		timeGetterMock           *mocks.TimeGetable
		ctx                      context.Context
		api                      *recurrentexpenses.CreateRecurrentExpenseImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		recurentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		timeGetterMock = &mocks.TimeGetable{}
		ctx = context.Background()
		api = recurrentexpenses.NewCreateRecurrentExpenseImpl(recurentExpensesRepoMock, expensesRepoMock, timeGetterMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		expensesRepoMock.AssertExpectations(T)
		recurentExpensesRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
	})

	It("creates recurrent expense with next month expense record", func() {
		var (
			expectedExpenseName        = faker.Name()
			expectedExpenseAmount      = faker.Latitude()
			expectedExpenseDescription = faker.Paragraph()
			expectedCreatedAt          = time.Date(2022, time.August, 1, 0, 0, 0, 0, time.Local)
			request                    = recurrentexpenses.CreateRecurrentExpenseInput{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: expectedExpenseDescription,
			}
			expectedRecurrentExpenseSaved = mongoentities.RecurrentExpense{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: expectedExpenseDescription,
			}
			expectedExpenseSaved = mongoentities.Expense{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: expectedExpenseDescription,
				IsRecurrent: true,
				CreatedAt:   &expectedCreatedAt,
			}
		)
		timeGetterMock.EXPECT().GetNextMonthAtFirtsDay().Return(expectedCreatedAt)
		recurentExpensesRepoMock.EXPECT().Save(ctx, &expectedRecurrentExpenseSaved).Return(nil)
		expensesRepoMock.EXPECT().Save(ctx, &expectedExpenseSaved).Return(nil)

		got, err := api.Create(ctx, &request)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.RecurrentExpense).To(Equal(&expectedRecurrentExpenseSaved))
		Expect(got.NextMonthExpense).To(Equal(&expectedExpenseSaved))
	})

})
