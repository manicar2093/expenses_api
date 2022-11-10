package mongorepos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/repos/mongorepos"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("IncomesRepo", func() {
	var (
		ctx  context.Context
		repo *mongorepos.IncomesMongoRepo
	)

	BeforeEach(func() {
		ctx = context.TODO()
		repo = mongorepos.NewIncomesMongoRepo(conn)

	})

	It("saves an mongoentities.Income in database", func() {

		var (
			expectedName        = faker.Name()
			expectedAmount      = faker.Latitude()
			expectedDescription = faker.Sentence()
			expectedIncome      = mongoentities.Income{
				Name:        expectedName,
				Amount:      expectedAmount,
				Description: expectedDescription,
			}
		)

		err := repo.Save(ctx, &expectedIncome)

		Expect(err).ToNot(HaveOccurred())
		Expect(expectedIncome.ID.String()).ToNot(BeEmpty())
		Expect(expectedIncome.Day).ToNot(BeZero())
		Expect(expectedIncome.Month).ToNot(BeZero())
		Expect(expectedIncome.Year).ToNot(BeZero())
		Expect(expectedIncome.CreatedAt).ToNot(BeZero())
		Expect(expectedIncome.UpdatedAt).To(BeNil())

		testfunc.DeleteOneByObjectID(ctx, conn.Collection("incomes"), expectedIncome.ID)
	})
})
