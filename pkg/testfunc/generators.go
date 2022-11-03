package testfunc

import "go.mongodb.org/mongo-driver/bson/primitive"

func GeneratePrimitiveObjectIDs(quantity uint) []primitive.ObjectID {
	var res []primitive.ObjectID
	for i := 0; i < int(quantity); i++ {
		res = append(res, primitive.NewObjectID())
	}
	return res
}
