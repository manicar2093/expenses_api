package connections_test

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/connections"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("Mongo", func() {
	var _ = Describe("MongoConnection", func() {
		var (
			ctx  context.Context
			db   *mongo.Database
			coll *mongo.Collection
		)

		BeforeEach(func() {
			ctx = context.TODO()
		})

		AfterEach(func() {
			db.Drop(ctx)
		})

		It("Should connect successfully", func() {
			db = connections.GetMongoConn()
			coll = db.Collection("testing_connection")
			res, err := coll.InsertOne(ctx, primitive.M{"is_success": true})

			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
		})

	})

})
