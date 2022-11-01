package testfunc

import (
	"context"

	"github.com/manicar2093/expenses_api/internal/repos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteManyByObjectID(ctx context.Context, coll *mongo.Collection, ids *mongo.InsertManyResult) {
	for _, v := range ids.InsertedIDs {
		coll.DeleteOne(ctx, bson.M{"_id": v.(primitive.ObjectID)}) //nolint: errcheck
	}
}

func DeleteManyByReposInsertManyResult(ctx context.Context, coll *mongo.Collection, ids *repos.InsertManyResult) {
	for _, v := range ids.InsertedIDs {
		coll.DeleteOne(ctx, bson.M{"_id": v}) //nolint: errcheck
	}
}

func DeleteOneByObjectID(ctx context.Context, coll *mongo.Collection, id primitive.ObjectID) {
	coll.DeleteOne(ctx, bson.M{"_id": id}) //nolint: errcheck
}
