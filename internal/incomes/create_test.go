package incomes_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/incomes"
	"github.com/manicar2093/expenses_api/mocks"
)

var _ = Describe("CreateImpl", func() {

	var (
		ctx             context.Context
		incomesRepoMock *mocks.IncomesRepository
		api             *incomes.CreateIncomeImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		incomesRepoMock = &mocks.IncomesRepository{}
		api = incomes.NewCreateIncomeImpl(incomesRepoMock)
	})

	AfterEach(func() {
		incomesRepoMock.AssertExpectations(GinkgoT())
	})

	It("create an entities.Incomes from schema", func() {
		var (
			expectedName        = faker.Name()
			expectedDescription = faker.Paragraph()
			expectedAmount      = faker.Latitude()
			incomeInput         = incomes.CreateIncomeInput{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
			}
			expextedIncomeEntity = mongoentities.Income{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
			}
		)
		incomesRepoMock.EXPECT().Save(ctx, &expextedIncomeEntity).Return(nil)

		got, err := api.Create(ctx, &incomeInput)

		Expect(err).ToNot(HaveOccurred())
		Expect(got).To(Equal(&expextedIncomeEntity))
	})
})
