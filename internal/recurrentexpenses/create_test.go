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

var _ = Describe("CreateRecurrentExpense", func() {

	var (
		expensesRepoMock         *mocks.ExpensesRepository
		recurentExpensesRepoMock *mocks.RecurrentExpenseRepo
		timeGetterMock           *mocks.TimeGetable
		ctx                      context.Context
		expectedUserID           string
		expectedUserIDAsUUID     uuid.UUID
		api                      *recurrentexpenses.RecurrentExpenseServiceImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		recurentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		timeGetterMock = &mocks.TimeGetable{}
		ctx = context.Background()
		expectedUserIDAsUUID = uuid.New()
		expectedUserID = expectedUserIDAsUUID.String()
		api = recurrentexpenses.NewCreateRecurrentExpense(recurentExpensesRepoMock, expensesRepoMock, timeGetterMock)
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
			expectedRecurrentExpenseID = uuid.New()
			expectedCreatedAt          = time.Date(2022, time.August, 1, 0, 0, 0, 0, time.Local)
			request                    = recurrentexpenses.CreateRecurrentExpenseInput{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: expectedExpenseDescription,
				UserID:      expectedUserID,
			}
			expectedRecurrentExpenseSaved = entities.RecurrentExpense{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: null.StringFrom(expectedExpenseDescription),
				UserID:      expectedUserIDAsUUID,
			}
			expectedRecurrentExpenseReturned = entities.RecurrentExpense{
				ID:          expectedRecurrentExpenseID,
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: null.StringFrom(expectedExpenseDescription),
				UserID:      expectedUserIDAsUUID,
			}
			expectedExpenseSaved = entities.Expense{
				Amount: expectedExpenseAmount,
				RecurrentExpenseID: uuid.NullUUID{
					UUID:  expectedRecurrentExpenseID,
					Valid: true,
				},
				CreatedAt: &expectedCreatedAt,
				UserID:    expectedUserIDAsUUID,
			}
		)
		timeGetterMock.EXPECT().GetNextMonthAtFirtsDay().Return(expectedCreatedAt)
		recurentExpensesRepoMock.EXPECT().Save(ctx, &expectedRecurrentExpenseSaved).Run(func(ctx context.Context, recExpense *entities.RecurrentExpense) {
			recExpense.ID = expectedRecurrentExpenseID
		}).Return(nil)
		expensesRepoMock.EXPECT().Save(ctx, &expectedExpenseSaved).Return(nil)

		got, err := api.CreateRecurrentExpense(ctx, &request)

		Expect(err).ToNot(HaveOccurred())
		Expect(got.RecurrentExpense).To(Equal(&expectedRecurrentExpenseReturned))
		Expect(got.NextMonthExpense).To(Equal(&expectedExpenseSaved))
	})
})
