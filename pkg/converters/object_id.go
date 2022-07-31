package converters

import "go.mongodb.org/mongo-driver/bson/primitive"

func TurnToObjectID(possibleID string) (primitive.ObjectID, error) {
	gotID, err := primitive.ObjectIDFromHex(possibleID)
	if err != nil {
		return primitive.ObjectID{}, &IDNotValidIDError{}
	}
	return gotID, nil
}
