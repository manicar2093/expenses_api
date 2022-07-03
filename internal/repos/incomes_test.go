package repos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/go-rel/reltest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
)

var _ = Describe("IncomesRepo", func() {
	var (
		ctx         context.Context
		relRepoMock *reltest.Repository
		repo        *repos.IncomesRepositoryImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		relRepoMock = reltest.New()
		repo = repos.NewIncomesRepositoryImpl(relRepoMock)

	})
	AfterEach(func() {
		T := GinkgoT()
		relRepoMock.AssertExpectations(T)
	})

	It("saves an entities.Income in database", func() {

		var (
			expectedName        = faker.Name()
			expectedAmount      = faker.Latitude()
			expectedDescription = faker.Sentence()
			expectedIncome      = entities.Income{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
			}
		)
		relRepoMock.ExpectTransaction(func(r *reltest.Repository) {
			r.ExpectInsert().For(&expectedIncome).Success()
		})

		err := repo.Save(ctx, &expectedIncome)

		Expect(err).ToNot(HaveOccurred())

	})
})
