package expenses_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/expenses"
	"github.com/manicar2093/expenses_api/mocks"
)

var _ = Describe("CreateImpl", func() {

	var (
		expenseRepoMock *mocks.ExpensesRepository
		ctx             context.Context
		api             *expenses.CreateExpenseImpl
	)

	BeforeEach(func() {
		expenseRepoMock = &mocks.ExpensesRepository{}
		ctx = context.TODO()
		api = expenses.NewCreateExpensesImpl(expenseRepoMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		expenseRepoMock.AssertExpectations(T)
	})

	It("creates a new expense from schema", func() {
		var (
			expectedName        = faker.Name()
			expectedDescription = faker.Paragraph()
			expectedAmount      = faker.Latitude()
			request             = expenses.CreateExpenseInput{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
			}
			expectedExpenseToSave = entities.Expense{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
				IsPaid:      true,
				IsRecurrent: false,
			}
		)
		expenseRepoMock.EXPECT().Save(ctx, &expectedExpenseToSave).Return(nil)

		got, err := api.Create(ctx, &request)

		Expect(err).ToNot(HaveOccurred())
		Expect(got).To(Equal(&expectedExpenseToSave))
	})

})
