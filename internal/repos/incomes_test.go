package repos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/manicar2093/expenses_api/internal/connections"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
)

var _ = Describe("IncomesRepo", func() {
	var (
		ctx  context.Context
		conn *mongo.Database
		repo *repos.IncomesRepositoryImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		conn = connections.GetMongoConn()
		repo = repos.NewIncomesRepositoryImpl(conn)

	})
	AfterEach(func() {
		conn.Drop(ctx)
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

		err := repo.Save(ctx, &expectedIncome)

		var incomeRegistered entities.Income
		conn.Collection("incomes").FindOne(ctx, bson.M{"_id": expectedIncome.ID}).Decode(&incomeRegistered)

		Expect(err).ToNot(HaveOccurred())
		Expect(expectedIncome.ID).To(BeAssignableToTypeOf(primitive.ObjectID{}))
		Expect(expectedIncome.CreatedAt).ToNot(BeZero())
		Expect(expectedIncome.CreatedAt).To(Equal(incomeRegistered.CreatedAt))
		Expect(expectedIncome.UpdatedAt).To(BeNil())

	})
})
