package converters

import "go.mongodb.org/mongo-driver/bson/primitive"

func TurnToObjectID(possibleID interface{}) (primitive.ObjectID, error) {
	var (
		emptyObjectID = primitive.ObjectID{}
		notValidError = &IDNotValidIDError{}
	)
	switch possibleID := possibleID.(type) {
	case string:
		gotID, err := primitive.ObjectIDFromHex(possibleID)
		if err != nil {
			return emptyObjectID, notValidError
		}
		return gotID, nil
	default:
		return emptyObjectID, notValidError
	}
}
