package connections

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoConn() *mongo.Database {
	var (
		ctx           = context.TODO()
		clientOptions = options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(os.Getenv("MONGO_DB"))
}
