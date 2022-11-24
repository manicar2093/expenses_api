package recurrentexpenses_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/recurrentexpenses"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	"github.com/manicar2093/expenses_api/pkg/validator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"
)

var _ = Describe("CreateRecurrentExpense", func() {

	var (
		expensesRepoMock         *mocks.ExpensesRepository
		recurentExpensesRepoMock *mocks.RecurrentExpenseRepo
		timeGetterMock           *mocks.TimeGetable
		validatorMock            *mocks.StructValidable
		ctx                      context.Context
		api                      *recurrentexpenses.RecurrentExpenseServiceImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		recurentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		timeGetterMock = &mocks.TimeGetable{}
		validatorMock = &mocks.StructValidable{}
		ctx = context.Background()
		api = recurrentexpenses.NewCreateRecurrentExpense(recurentExpensesRepoMock, expensesRepoMock, timeGetterMock, validatorMock)
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
			}
			expectedRecurrentExpenseSaved = entities.RecurrentExpense{
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: null.StringFrom(expectedExpenseDescription),
			}
			expectedRecurrentExpenseReturned = entities.RecurrentExpense{
				ID:          expectedRecurrentExpenseID,
				Name:        expectedExpenseName,
				Amount:      expectedExpenseAmount,
				Description: null.StringFrom(expectedExpenseDescription),
			}
			expectedExpenseSaved = entities.Expense{
				Amount: expectedExpenseAmount,
				RecurrentExpenseID: uuid.NullUUID{
					UUID:  expectedRecurrentExpenseID,
					Valid: true,
				},
				CreatedAt: &expectedCreatedAt,
			}
		)
		validatorMock.EXPECT().ValidateStruct(&request).Return(nil)
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

	When("request is not valid", Label(testfunc.IntegrationTest), func() {
		It("return an error", func() {
			var invalidRequest = recurrentexpenses.CreateRecurrentExpenseInput{}

			integrationTestApi := recurrentexpenses.NewCreateRecurrentExpense(recurentExpensesRepoMock, expensesRepoMock, timeGetterMock, validator.NewGooKitValidator())

			got, err := integrationTestApi.CreateRecurrentExpense(ctx, &invalidRequest)

			Expect(got).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})
})
