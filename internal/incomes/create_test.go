package incomes_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/incomes"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/goption"
)

var _ = Describe("CreateImpl", func() {

	var (
		ctx             context.Context
		incomesRepoMock *mocks.IncomesRepository
		userID          goption.Optional[uuid.UUID]
		api             *incomes.IncomeServiceImpl
	)

	BeforeEach(func() {
		T := GinkgoT()
		ctx = context.TODO()
		incomesRepoMock = mocks.NewIncomesRepository(T)
		userID = goption.Of(uuid.New())
		api = incomes.NewIncomeServiceImpl(incomesRepoMock)
	})

	It("create an entities.Incomes from schema", func() {
		var (
			expectedName         = faker.Name()
			expectedDescription  = faker.Paragraph()
			expectedAmount       = faker.Latitude()
			expextedIncomeEntity = entities.Income{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
				UserID:      userID,
			}
			incomeInput = incomes.CreateIncomeInput{
				expextedIncomeEntity,
			}
		)
		incomesRepoMock.EXPECT().Save(ctx, &expextedIncomeEntity).Return(nil)

		got, err := api.Create(ctx, &incomeInput)

		Expect(err).ToNot(HaveOccurred())
		Expect(got).To(Equal(&expextedIncomeEntity))
	})
})
