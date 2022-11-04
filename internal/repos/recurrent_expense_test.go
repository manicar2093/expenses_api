package repos_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/dates"
	"github.com/manicar2093/expenses_api/pkg/periodtypes"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
)

var _ = Describe("RecurrentExpense", func() {
	var (
		ctx        context.Context
		coll       *mongo.Collection
		timeGetter *mocks.TimeGetable
		repo       *repos.RecurrentExpenseRepoImpl
	)

	BeforeEach(func() {
		ctx = context.TODO()
		coll = conn.Collection("recurrent_expenses")
		timeGetter = &mocks.TimeGetable{}
		repo = repos.NewRecurrentExpenseRepoImpl(conn, timeGetter)

	})

	Describe("Save", func() {
		It("saves an instance", func() {
			var (
				expectedLastCreationDate = time.Now()
				toSave                   = entities.RecurrentExpense{
					Name:             faker.Name(),
					Amount:           faker.Latitude(),
					Description:      faker.Paragraph(),
					Periodicity:      periodtypes.Monthly,
					LastCreationDate: &expectedLastCreationDate,
				}
			)

			err := repo.Save(ctx, &toSave)

			Expect(err).ToNot(HaveOccurred())
			Expect(toSave.ID.String()).ToNot(BeEmpty())
			Expect(toSave.Periodicity).To(Equal(periodtypes.Monthly))
			Expect(toSave.LastCreationDate).To(Equal(&expectedLastCreationDate))
			Expect(toSave.CreatedAt).ToNot(BeZero())
			Expect(toSave.UpdatedAt).To(BeNil())

			testfunc.DeleteOneByObjectID(ctx, coll, toSave.ID)
		})

		When("recurrent expense exists", func() {
			It("returns an AlreadyExists", func() {
				var (
					expectedName = "testing"
					Saved        = entities.RecurrentExpense{
						Name:        expectedName,
						Amount:      faker.Latitude(),
						Description: faker.Paragraph(),
					}
					toSave = entities.RecurrentExpense{
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
			Expect(got).To(HaveLen(len(dataSaved)))

			testfunc.DeleteManyByObjectID(ctx, coll, inserted)
		})
	})

	Describe("Update", func() {
		It("changes recurrent expense data in db", func() {
			var (
				toUpdate = entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.Monthly,
				}
				expectedDescription      = faker.Paragraph()
				expectedLastCreationDate = dates.NormalizeDate(time.Now())
				expectedToday            = dates.NormalizeDate(time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local))
			)
			timeGetter.EXPECT().GetCurrentTime().Return(expectedToday)
			coll.InsertOne(ctx, &toUpdate) //nolint:errcheck
			toUpdate.Description = expectedDescription
			toUpdate.LastCreationDate = &expectedLastCreationDate

			err := repo.Update(ctx, &toUpdate)
			var updated entities.RecurrentExpense
			coll.FindOne(ctx, primitive.D{{Key: "_id", Value: toUpdate.ID}}).Decode(&updated) //nolint:errcheck

			Expect(err).ToNot(HaveOccurred())
			Expect(updated.ID).To(Equal(toUpdate.ID))
			Expect(updated.Description).To(Equal(expectedDescription))
			Expect(updated.Periodicity).To(Equal(toUpdate.Periodicity))
			Expect(updated.LastCreationDate).To(Equal(&expectedLastCreationDate))
			Expect(updated.CreatedAt).To(Equal(toUpdate.CreatedAt))
			Expect(updated.UpdatedAt).To(Equal(&expectedToday))

			testfunc.DeleteOneByObjectID(ctx, coll, toUpdate.ID)
		})
	})
})
