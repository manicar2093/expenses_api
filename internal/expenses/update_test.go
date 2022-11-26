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
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	"github.com/manicar2093/expenses_api/pkg/validator"
)

var _ = Describe("Update", func() {

	var (
		expensesRepoMock *mocks.ExpensesRepository
		validatorMock    *mocks.StructValidable
		timeGettable     *mocks.TimeGetable
		ctx              context.Context
		api              *expenses.ExpenseServiceImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		validatorMock = &mocks.StructValidable{}
		timeGettable = &mocks.TimeGetable{}
		ctx = context.Background()
		api = expenses.NewExpenseServiceImpl(expensesRepoMock, timeGettable, validatorMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		expensesRepoMock.AssertExpectations(T)
		validatorMock.AssertExpectations(T)
		timeGettable.AssertExpectations(T)
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
		validatorMock.EXPECT().ValidateStruct(&expectedUpdateInput).Return(nil)
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
			validatorMock.EXPECT().ValidateStruct(&expectedUpdateInput).Return(nil)
			expensesRepoMock.EXPECT().Update(ctx, &expectedRepoCall).Return(nil)

			err := api.UpdateExpense(ctx, &expectedUpdateInput)

			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("request is not valid", Label(testfunc.IntegrationTest), func() {
		It("return an error", func() {
			var invalidRequest = expenses.UpdateExpenseInput{
				ID: faker.Name(),
			}

			integrationTestApi := expenses.NewExpenseServiceImpl(expensesRepoMock, timeGettable, validator.NewGooKitValidator())

			err := integrationTestApi.UpdateExpense(ctx, &invalidRequest)

			Expect(err).To(HaveOccurred())
		})
	})

})
