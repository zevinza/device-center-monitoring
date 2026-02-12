package lib

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HexToObjectID(hex string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(strings.TrimSpace(hex))
}

func ObjectIDToHex(objectID primitive.ObjectID) string {
	return objectID.Hex()
}
