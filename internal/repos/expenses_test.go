package repos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/go-rel/reltest"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExpensesImpl", func() {
	var (
		ctx         context.Context
		relRepoMock *reltest.Repository
		repo        *repos.ExpensesRepositoryImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		relRepoMock = reltest.New()
		repo = repos.NewExpensesRepositoryImpl(relRepoMock)

	})
	AfterEach(func() {
		T := GinkgoT()
		relRepoMock.AssertExpectations(T)
	})

	It("saves an entities.Expense in database", func() {

		var (
			expectedName        = faker.Name()
			expectedAmount      = faker.Latitude()
			expectedDescription = faker.Sentence()
			expectedExpense     = entities.Expense{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
			}
		)
		relRepoMock.ExpectTransaction(func(r *reltest.Repository) {
			r.ExpectInsert().For(&expectedExpense).Success()
		})

		err := repo.Save(ctx, &expectedExpense)

		Expect(err).ToNot(HaveOccurred())

	})
})
