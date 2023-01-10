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
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	"github.com/manicar2093/expenses_api/pkg/validator"
	"github.com/manicar2093/goption"
)

var _ = Describe("CreateImpl", func() {

	var (
		ctx             context.Context
		incomesRepoMock *mocks.IncomesRepository
		validatorMock   *mocks.StructValidable
		userID          goption.Optional[uuid.UUID]
		api             *incomes.IncomeServiceImpl
	)

	BeforeEach(func() {
		T := GinkgoT()
		ctx = context.TODO()
		incomesRepoMock = mocks.NewIncomesRepository(T)
		validatorMock = mocks.NewStructValidable(T)
		userID = goption.Of(uuid.New())
		api = incomes.NewIncomeServiceImpl(incomesRepoMock, validatorMock)
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
		validatorMock.EXPECT().ValidateStruct(&incomeInput).Return(nil)
		incomesRepoMock.EXPECT().Save(ctx, &expextedIncomeEntity).Return(nil)

		got, err := api.Create(ctx, &incomeInput)

		Expect(err).ToNot(HaveOccurred())
		Expect(got).To(Equal(&expextedIncomeEntity))
	})

	When("request is not valid", Label(testfunc.IntegrationTest), func() {
		It("return an error", func() {
			var invalidRequest = incomes.CreateIncomeInput{}

			integrationTestApi := incomes.NewIncomeServiceImpl(incomesRepoMock, validator.NewGooKitValidator())

			got, err := integrationTestApi.Create(ctx, &invalidRequest)

			Expect(got).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})
})
