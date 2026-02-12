package sensorreadingdomain

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorreadingrepo"
	"api/app/master-service/repository/sensorrepo"
	"api/constant"
	"api/services/queue"
	"api/utils/resp"
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorReadingDomain interface {
	Ingest(ctx context.Context, req *model.SensorIngestRequest) error
}

type sensorReadingDomain struct {
	readingRepo sensorreadingrepo.SensorReadingRepository
	deviceRepo  devicerepo.DeviceRepository
	sensorRepo  sensorrepo.SensorRepository
	queue       *queue.RedisQueue
}

func New(
	readingRepo sensorreadingrepo.SensorReadingRepository,
	deviceRepo devicerepo.DeviceRepository,
	sensorRepo sensorrepo.SensorRepository,
	queue *queue.RedisQueue,
) SensorReadingDomain {
	return &sensorReadingDomain{
		readingRepo: readingRepo,
		deviceRepo:  deviceRepo,
		sensorRepo:  sensorRepo,
		queue:       queue,
	}
}

func (d *sensorReadingDomain) Ingest(ctx context.Context, req *model.SensorIngestRequest) error {
	oid, err := primitive.ObjectIDFromHex(req.SensorID)
	if err != nil {
		return resp.ErrorBadRequest("invalid sensor id")
	}

	sensor, err := d.sensorRepo.GetByID(ctx, oid)
	if err != nil {
		return resp.ErrorNotFound("sensor not found")
	}

	device, err := d.deviceRepo.GetByID(ctx, sensor.DeviceID)
	if err != nil {
		return resp.ErrorNotFound("device not found")
	}

	ts := time.Now().UTC()
	if req.Timestamp != nil {
		ts = req.Timestamp.UTC()
	}

	// Create sensor reading
	reading := &model.SensorReading{
		DeviceID:   device.ID,
		SensorID:   sensor.ID,
		Value:      req.Value,
		Timestamp:  ts,
		Status:     constant.SensorReadingStatus_Pending,
		RetryCount: 0,
		CreatedAt:  ts,
		UpdatedAt:  ts,
	}

	if err := d.readingRepo.Create(ctx, reading); err != nil {
		return resp.ErrorInternal(err.Error())
	}

	msg := &model.QueueMessage{
		ReadingID: reading.ID.Hex(),
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return resp.ErrorInternal(err.Error())
	}

	// push to queue
	if err := d.queue.Enqueue(ctx, jsonData); err != nil {
		return resp.ErrorInternal(err.Error())
	}

	return nil
}
