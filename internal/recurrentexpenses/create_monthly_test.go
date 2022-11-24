package recurrentexpenses_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"

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
		api                       *recurrentexpenses.RecurrentExpenseServiceImpl

		expectedGetNextMonthAtFirtsDayReturn time.Time
		expectedMonth                        uint
		expectedName1                        string
		expectedName2                        string
		expectedName3                        string
		expectedName4                        string
		expectedAmount1                      float64
		expectedAmount2                      float64
		expectedAmount3                      float64
		expectedAmount4                      float64
		expectedRecurrenteExpensesFound      []*entities.RecurrentExpense
		expectedExpensesToCreate             []*entities.Expense
	)

	BeforeEach(func() {
		recurrentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		expensesRepoMock = &mocks.ExpensesRepository{}
		timeGetableMock = &mocks.TimeGetable{}
		ctx = context.Background()
		api = recurrentexpenses.NewCreateMonthlyRecurrentExpensesImpl(recurrentExpensesRepoMock, expensesRepoMock, timeGetableMock)

		expectedGetNextMonthAtFirtsDayReturn = time.Date(2022, time.September, 1, 0, 0, 0, 0, time.Local)
		expectedMonth = uint(expectedGetNextMonthAtFirtsDayReturn.Month())
		expectedName1 = faker.Name()
		expectedName2 = faker.Name()
		expectedName3 = faker.Name()
		expectedName4 = faker.Name()
		expectedAmount1 = faker.Latitude()
		expectedAmount2 = faker.Latitude()
		expectedAmount3 = faker.Latitude()
		expectedAmount4 = faker.Latitude()
		expectedRecurrenteExpensesFound = []*entities.RecurrentExpense{
			{Name: expectedName1, Amount: expectedAmount1},
			{Name: expectedName2, Amount: expectedAmount2},
			{Name: expectedName3, Amount: expectedAmount3},
			{Name: expectedName4, Amount: expectedAmount4},
		}
		expectedExpensesToCreate = []*entities.Expense{
			{Name: null.StringFrom(expectedName1), Amount: expectedAmount1, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
			{Name: null.StringFrom(expectedName2), Amount: expectedAmount2, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
			{Name: null.StringFrom(expectedName3), Amount: expectedAmount3, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
			{Name: null.StringFrom(expectedName4), Amount: expectedAmount4, CreatedAt: &expectedGetNextMonthAtFirtsDayReturn},
		}
	})

	AfterEach(func() {
		T := GinkgoT()
		recurrentExpensesRepoMock.AssertExpectations(T)
		expensesRepoMock.AssertExpectations(T)
		timeGetableMock.AssertExpectations(T)
	})

	It("finds all recurrent expenses and creates them as expenses", func() {
		timeGetableMock.EXPECT().GetNextMonthAtFirtsDay().Return(expectedGetNextMonthAtFirtsDayReturn)
		recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(expectedRecurrenteExpensesFound, nil)
		expensesRepoMock.EXPECT().FindByNameAndMonthAndIsRecurrent(ctx, expectedMonth, expectedRecurrenteExpensesFound[0].Name).Return(nil, &repos.NotFoundError{}).Once()
		expensesRepoMock.EXPECT().FindByNameAndMonthAndIsRecurrent(ctx, expectedMonth, expectedRecurrenteExpensesFound[1].Name).Return(nil, &repos.NotFoundError{}).Once()
		expensesRepoMock.EXPECT().FindByNameAndMonthAndIsRecurrent(ctx, expectedMonth, expectedRecurrenteExpensesFound[2].Name).Return(nil, &repos.NotFoundError{}).Once()
		expensesRepoMock.EXPECT().FindByNameAndMonthAndIsRecurrent(ctx, expectedMonth, expectedRecurrenteExpensesFound[3].Name).Return(nil, &repos.NotFoundError{}).Once()
		expensesRepoMock.EXPECT().Save(ctx, expectedExpensesToCreate[0]).Return(nil).Once()
		expensesRepoMock.EXPECT().Save(ctx, expectedExpensesToCreate[1]).Return(nil).Once()
		expensesRepoMock.EXPECT().Save(ctx, expectedExpensesToCreate[2]).Return(nil).Once()
		expensesRepoMock.EXPECT().Save(ctx, expectedExpensesToCreate[3]).Return(nil).Once()

		got, err := api.CreateMonthlyRecurrentExpenses(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.ExpensesCreated).To(Equal(expectedExpensesToCreate))
	})

	When("some expenses was already created as recurrent expense", func() {
		It("avoids creation", func() {
			timeGetableMock.EXPECT().GetNextMonthAtFirtsDay().Return(expectedGetNextMonthAtFirtsDayReturn)
			recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(expectedRecurrenteExpensesFound, nil)
			expensesRepoMock.EXPECT().FindByNameAndMonthAndIsRecurrent(ctx, expectedMonth, expectedRecurrenteExpensesFound[0].Name).Return(expectedExpensesToCreate[0], nil).Once()
			expensesRepoMock.EXPECT().FindByNameAndMonthAndIsRecurrent(ctx, expectedMonth, expectedRecurrenteExpensesFound[1].Name).Return(expectedExpensesToCreate[1], nil).Once()
			expensesRepoMock.EXPECT().FindByNameAndMonthAndIsRecurrent(ctx, expectedMonth, expectedRecurrenteExpensesFound[2].Name).Return(nil, &repos.NotFoundError{}).Once()
			expensesRepoMock.EXPECT().FindByNameAndMonthAndIsRecurrent(ctx, expectedMonth, expectedRecurrenteExpensesFound[3].Name).Return(nil, &repos.NotFoundError{}).Once()
			expensesRepoMock.EXPECT().Save(ctx, expectedExpensesToCreate[2]).Return(nil).Once()
			expensesRepoMock.EXPECT().Save(ctx, expectedExpensesToCreate[3]).Return(nil).Once()

			_, err := api.CreateMonthlyRecurrentExpenses(ctx)

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
