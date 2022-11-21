package recurrentexpenses_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"
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
			expectedRecurrentExpenseSaved = entities.RecurrentExpense{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: null.StringFrom(expectedExpenseDescription),
			}
			expectedExpenseSaved = entities.Expense{
				Name:   expectedExpenseName,
				Amount: expectedExpenseAmount,
				RecurrentExpenseID: uuid.NullUUID{
					UUID:  expectedRecurrentExpenseSaved.ID,
					Valid: true,
				},
				Description: null.StringFrom(expectedExpenseDescription),
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
