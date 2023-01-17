package expenses_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/mocks"
)

var _ = Describe("Update", func() {

	var (
		expensesRepoMock *mocks.ExpensesRepository
		timeGettable     *mocks.TimeGetable
		ctx              context.Context
		api              *expenses.ExpenseServiceImpl
	)

	BeforeEach(func() {
		T := GinkgoT()
		expensesRepoMock = mocks.NewExpensesRepository(T)
		timeGettable = mocks.NewTimeGetable(T)
		ctx = context.Background()
		api = expenses.NewExpenseServiceImpl(expensesRepoMock, timeGettable)
	})

	It("calls repo to change expense stored data", func() {
		var (
			expectedExpenseID          = uuid.New().String()
			expectedExpenseIDAsUUID    = uuid.MustParse(expectedExpenseID)
			expectedExpenseName        = faker.Name()
			expectedExpenseAmount      = faker.Latitude()
			expectedExpenseDescription = faker.Paragraph()
			expectedUpdateInput        = expenses.UpdateExpenseInput{
				ID:          expectedExpenseID,
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: expectedExpenseDescription,
			}
			expectedExpenseFindByIDReturn = entities.Expense{
				ID:                 expectedExpenseIDAsUUID,
				Name:               null.StringFrom(faker.Name()),
				Amount:             faker.Latitude(),
				Description:        null.StringFrom(faker.Paragraph()),
				RecurrentExpenseID: uuid.NullUUID{},
				Day:                1,
				Month:              2,
				Year:               2022,
				IsPaid:             true,
			}
			expectedRepoCall = repos.UpdateExpenseInput{
				ID:          uuid.MustParse(expectedExpenseID),
				Name:        null.StringFrom(expectedExpenseName),
				Amount:      expectedExpenseAmount,
				Description: null.StringFrom(expectedExpenseDescription),
			}
		)
		expensesRepoMock.EXPECT().FindByID(ctx, expectedExpenseIDAsUUID).Return(&expectedExpenseFindByIDReturn, nil)
		expensesRepoMock.EXPECT().Update(ctx, &expectedRepoCall).Return(nil)

		err := api.UpdateExpense(ctx, &expectedUpdateInput)

		Expect(err).ToNot(HaveOccurred())
	})

	When("expense is recurrent", func() {
		It("just update amount", func() {
			var (
				expectedExpenseID          = uuid.New().String()
				expectedExpenseIDAsUUID    = uuid.MustParse(expectedExpenseID)
				expectedExpenseName        = faker.Name()
				expectedExpenseAmount      = faker.Latitude()
				expectedExpenseDescription = faker.Paragraph()
				expectedUpdateInput        = expenses.UpdateExpenseInput{
					ID:          expectedExpenseID,
					Name:        expectedExpenseName,
					Amount:      expectedExpenseAmount,
					Description: expectedExpenseDescription,
				}
				expectedExpenseFindByIDReturn = entities.Expense{
					ID:          expectedExpenseIDAsUUID,
					Name:        null.StringFrom(faker.Name()),
					Amount:      faker.Latitude(),
					Description: null.StringFrom(faker.Paragraph()),
					RecurrentExpenseID: uuid.NullUUID{
						UUID:  uuid.New(),
						Valid: true,
					},
					Day:    1,
					Month:  2,
					Year:   2022,
					IsPaid: true,
				}
				expectedRepoCall = repos.UpdateExpenseInput{
					ID:     uuid.MustParse(expectedExpenseID),
					Amount: expectedExpenseAmount,
				}
			)
			expensesRepoMock.EXPECT().FindByID(ctx, expectedExpenseIDAsUUID).Return(&expectedExpenseFindByIDReturn, nil)
			expensesRepoMock.EXPECT().Update(ctx, &expectedRepoCall).Return(nil)

			err := api.UpdateExpense(ctx, &expectedUpdateInput)

			Expect(err).ToNot(HaveOccurred())
		})
	})

})
