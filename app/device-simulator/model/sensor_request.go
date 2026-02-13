package model

import "time"

type SensorIngestRequest struct {
	SensorID  string     `json:"sensor_id" validate:"required"`
	Value     float64    `json:"value" validate:"required"`
	Timestamp *time.Time `json:"timestamp"`
}
