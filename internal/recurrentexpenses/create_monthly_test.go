package recurrentexpenses_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/mocks"
)

var _ = Describe("CreateMonthly", func() {

	var (
		recurrentExpensesRepoMock *mocks.RecurrentExpenseRepo
		expensesRepoMock          *mocks.ExpensesRepository
		timeGetableMock           *mocks.TimeGetable
		ctx                       context.Context
		api                       *recurrentexpenses.CreateMonthlyRecurrentExpensesImpl
	)

	BeforeEach(func() {
		recurrentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		expensesRepoMock = &mocks.ExpensesRepository{}
		timeGetableMock = &mocks.TimeGetable{}
		ctx = context.Background()
		api = recurrentexpenses.NewCreateMonthlyRecurrentExpensesImpl(recurrentExpensesRepoMock, expensesRepoMock, timeGetableMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		recurrentExpensesRepoMock.AssertExpectations(T)
		expensesRepoMock.AssertExpectations(T)
		timeGetableMock.AssertExpectations(T)
	})

	It("finds all recurrent expenses and creates them as expenses", func() {
		var (
			expectedGetNextMonthAtFirtsDayReturn = time.Date(2022, time.September, 1, 0, 0, 0, 0, time.Local)
			expectedName1                        = faker.Name()
			expectedName2                        = faker.Name()
			expectedName3                        = faker.Name()
			expectedName4                        = faker.Name()
			expectedAmount1                      = faker.Latitude()
			expectedAmount2                      = faker.Latitude()
			expectedAmount3                      = faker.Latitude()
			expectedAmount4                      = faker.Latitude()
			expectedRecurrenteExpensesFound      = []entities.RecurrentExpense{
				{Name: expectedName1, Amount: expectedAmount1},
				{Name: expectedName2, Amount: expectedAmount2},
				{Name: expectedName3, Amount: expectedAmount3},
				{Name: expectedName4, Amount: expectedAmount4},
			}
			expectedExpensesToCreate = []entities.Expense{
				{Name: expectedName1, Amount: expectedAmount1, IsRecurrent: true, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
				{Name: expectedName2, Amount: expectedAmount2, IsRecurrent: true, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
				{Name: expectedName3, Amount: expectedAmount3, IsRecurrent: true, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
				{Name: expectedName4, Amount: expectedAmount4, IsRecurrent: true, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
			}
		)
		timeGetableMock.EXPECT().GetNextMonthAtFirtsDay().Return(expectedGetNextMonthAtFirtsDayReturn)
		recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(&expectedRecurrenteExpensesFound, nil)
		expensesRepoMock.EXPECT().FindByNameAndIsRecurrent(ctx, expectedRecurrenteExpensesFound[0].Name).Return(nil, &repos.NotFoundError{}).Once()
		expensesRepoMock.EXPECT().FindByNameAndIsRecurrent(ctx, expectedRecurrenteExpensesFound[1].Name).Return(nil, &repos.NotFoundError{}).Once()
		expensesRepoMock.EXPECT().FindByNameAndIsRecurrent(ctx, expectedRecurrenteExpensesFound[2].Name).Return(nil, &repos.NotFoundError{}).Once()
		expensesRepoMock.EXPECT().FindByNameAndIsRecurrent(ctx, expectedRecurrenteExpensesFound[3].Name).Return(nil, &repos.NotFoundError{}).Once()
		expensesRepoMock.EXPECT().Save(ctx, &expectedExpensesToCreate[0]).Return(nil).Once()
		expensesRepoMock.EXPECT().Save(ctx, &expectedExpensesToCreate[1]).Return(nil).Once()
		expensesRepoMock.EXPECT().Save(ctx, &expectedExpensesToCreate[2]).Return(nil).Once()
		expensesRepoMock.EXPECT().Save(ctx, &expectedExpensesToCreate[3]).Return(nil).Once()

		err := api.CreateMonthlyRecurrentExpenses(ctx)

		Expect(err).ToNot(HaveOccurred())
	})

	When("some expenses was already created as recurrent expense", func() {
		It("avoids creation", func() {
			var (
				expectedGetNextMonthAtFirtsDayReturn = time.Date(2022, time.September, 1, 0, 0, 0, 0, time.Local)
				expectedName1                        = faker.Name()
				expectedName2                        = faker.Name()
				expectedName3                        = faker.Name()
				expectedName4                        = faker.Name()
				expectedAmount1                      = faker.Latitude()
				expectedAmount2                      = faker.Latitude()
				expectedAmount3                      = faker.Latitude()
				expectedAmount4                      = faker.Latitude()
				expectedRecurrenteExpensesFound      = []entities.RecurrentExpense{
					{Name: expectedName1, Amount: expectedAmount1},
					{Name: expectedName2, Amount: expectedAmount2},
					{Name: expectedName3, Amount: expectedAmount3},
					{Name: expectedName4, Amount: expectedAmount4},
				}
				expectedExpensesToCreate = []entities.Expense{
					{Name: expectedName1, Amount: expectedAmount1, IsRecurrent: true, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
					{Name: expectedName2, Amount: expectedAmount2, IsRecurrent: true, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
					{Name: expectedName3, Amount: expectedAmount3, IsRecurrent: true, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
					{Name: expectedName4, Amount: expectedAmount4, IsRecurrent: true, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
				}
			)
			timeGetableMock.EXPECT().GetNextMonthAtFirtsDay().Return(expectedGetNextMonthAtFirtsDayReturn)
			recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(&expectedRecurrenteExpensesFound, nil)
			expensesRepoMock.EXPECT().FindByNameAndIsRecurrent(ctx, expectedRecurrenteExpensesFound[0].Name).Return(&expectedExpensesToCreate[0], nil).Once()
			expensesRepoMock.EXPECT().FindByNameAndIsRecurrent(ctx, expectedRecurrenteExpensesFound[1].Name).Return(&expectedExpensesToCreate[1], nil).Once()
			expensesRepoMock.EXPECT().FindByNameAndIsRecurrent(ctx, expectedRecurrenteExpensesFound[2].Name).Return(nil, &repos.NotFoundError{}).Once()
			expensesRepoMock.EXPECT().FindByNameAndIsRecurrent(ctx, expectedRecurrenteExpensesFound[3].Name).Return(nil, &repos.NotFoundError{}).Once()
			expensesRepoMock.EXPECT().Save(ctx, &expectedExpensesToCreate[2]).Return(nil).Once()
			expensesRepoMock.EXPECT().Save(ctx, &expectedExpensesToCreate[3]).Return(nil).Once()

			err := api.CreateMonthlyRecurrentExpenses(ctx)

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
