package model

import "time"

type SensorIngestRequest struct {
	SensorID  string     `json:"sensor_id" validate:"required"`
	Value     string     `json:"value" validate:"required"`
	Timestamp *time.Time `json:"timestamp"`
}
