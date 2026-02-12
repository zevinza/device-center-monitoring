package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sensor is a managed sensor attached to a device (one device -> many sensors).
type Sensor struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeviceID  primitive.ObjectID `bson:"device_id" json:"device_id"`
	Name      string             `bson:"name" json:"name"`
	Unit      string             `bson:"unit" json:"unit"`
	Type      string             `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type SensorCreateRequest struct {
	DeviceID primitive.ObjectID `json:"device_id"`
	Name     string             `json:"name"`
	Unit     string             `json:"unit"`
	Type     string             `json:"type"`
}

type SensorUpdateRequest struct {
	Name *string `json:"name,omitempty"`
	Unit *string `json:"unit,omitempty"`
	Type *string `json:"type,omitempty"`
}
