package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SensorReading is the time-series data received from devices.
type SensorReading struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	SensorID   *uuid.UUID         `bson:"sensor_id" json:"sensor_id"`
	Value      any                `bson:"value" json:"value"`
	Timestamp  *time.Time          `bson:"timestamp" json:"timestamp"`
	Status     string             `bson:"status" json:"status"`
	RetryCount int                `bson:"retry_count" json:"retry_count"`
}

type SensorIngestRequest struct {
	SensorID  *uuid.UUID `json:"sensor_id" validate:"required"`
	Value     any        `json:"value" validate:"required"`
	Timestamp *time.Time `json:"timestamp"`
}
