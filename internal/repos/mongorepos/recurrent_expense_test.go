package mongorepos_test

import (
	"context"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/manicar2093/expenses_api/internal/entities/mongoentities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/internal/repos/mongorepos"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("RecurrentExpense", func() {
	var (
		ctx  context.Context
		coll *mongo.Collection
		repo *mongorepos.RecurrentExpenseMongoRepo
	)

	BeforeEach(func() {
		ctx = context.TODO()
		coll = conn.Collection("recurrent_expenses")
		repo = mongorepos.NewRecurrentExpenseMongoRepo(conn)

	})

	Describe("Save", func() {
		It("saves an instance", func() {
			var (
				toSave = mongoentities.RecurrentExpense{
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
				}
			)

			err := repo.Save(ctx, &toSave)

			Expect(err).ToNot(HaveOccurred())
			Expect(toSave.ID.String()).ToNot(BeEmpty())
			Expect(toSave.CreatedAt).ToNot(BeZero())
			Expect(toSave.UpdatedAt).To(BeNil())

			testfunc.DeleteOneByObjectID(ctx, coll, toSave.ID)
		})

		When("recurrent expense exists", func() {
			It("returns an AlreadyExists", func() {
				var (
					expectedName = "testing"
					Saved        = mongoentities.RecurrentExpense{
						Name:        expectedName,
						Amount:      faker.Latitude(),
						Description: faker.Paragraph(),
					}
					toSave = mongoentities.RecurrentExpense{
						Name:        expectedName,
						Amount:      faker.Latitude(),
						Description: faker.Paragraph(),
					}
				)
				repo.Save(ctx, &Saved) //nolint: errcheck

				err := repo.Save(ctx, &toSave)

				Expect(err).To(BeAssignableToTypeOf(&repos.AlreadyExistsError{}))
				Expect(err.(*repos.AlreadyExistsError).Entity).To(Equal("RecurrentExpense"))
				Expect(err.(*repos.AlreadyExistsError).Identifier).To(Equal(Saved.Name))

				testfunc.DeleteOneByObjectID(ctx, coll, Saved.ID)
				testfunc.DeleteOneByObjectID(ctx, coll, toSave.ID)
			})
		})
	})

	Describe("FindByName", func() {
		It("returns a pointer of found data", func() {
			var (
				expectedName = faker.Name()
				dataSaved    = []interface{}{
					bson.D{{Key: "name", Value: expectedName}},
					bson.D{{Key: "name", Value: faker.Name()}},
					bson.D{{Key: "name", Value: faker.Name()}},
				}
			)
			inserted, _ := coll.InsertMany(ctx, dataSaved)

			got, err := repo.FindByName(ctx, expectedName)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.Name).To(Equal(expectedName))

			testfunc.DeleteManyByObjectID(ctx, coll, inserted)
		})
	})

	Describe("FindAll", func() {
		It("gets all registered recurrent expenses", func() {
			var (
				dataSaved = []interface{}{
					bson.D{{Key: "name", Value: faker.Name()}},
					bson.D{{Key: "name", Value: faker.Name()}},
					bson.D{{Key: "name", Value: faker.Name()}},
					bson.D{{Key: "name", Value: faker.Name()}},
					bson.D{{Key: "name", Value: faker.Name()}},
				}
			)
			inserted, _ := coll.InsertMany(ctx, dataSaved)

			got, err := repo.FindAll(ctx)

			Expect(err).ToNot(HaveOccurred())
			Expect(*got).To(HaveLen(len(dataSaved)))

			testfunc.DeleteManyByObjectID(ctx, coll, inserted)
		})
	})
})
