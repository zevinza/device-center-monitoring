package sensorreadingdomain

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorreadingrepo"
	"api/app/master-service/repository/sensorrepo"
	"api/constant"
	"api/lib"
	"api/services/queue"
	"api/utils/resp"
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

type SensorReadingDomain interface {
	Ingest(ctx context.Context, req *model.SensorIngestRequest) error
}

type sensorReadingDomain struct {
	readingRepo sensorreadingrepo.SensorReadingRepository
	deviceRepo  devicerepo.DeviceRepository
	sensorRepo  sensorrepo.SensorRepository
	db          *gorm.DB
	queue       *queue.RedisQueue
}

func New(
	db *gorm.DB,
	queue *queue.RedisQueue,
	readingRepo sensorreadingrepo.SensorReadingRepository,
	deviceRepo devicerepo.DeviceRepository,
	sensorRepo sensorrepo.SensorRepository,
) SensorReadingDomain {
	return &sensorReadingDomain{
		readingRepo: readingRepo,
		deviceRepo:  deviceRepo,
		sensorRepo:  sensorRepo,
		queue:       queue,
		db:          db,
	}
}

func (d *sensorReadingDomain) Ingest(ctx context.Context, req *model.SensorIngestRequest) error {
	sensor, err := d.sensorRepo.GetByID(ctx, d.db, req.SensorID)
	if err != nil {
		return resp.ErrorNotFound("sensor not found")
	}

	// Create sensor reading
	reading := &model.SensorReading{
		SensorID:   lib.StringUUID(sensor.ID),
		Value:      req.Value,
		Timestamp:  req.Timestamp,
		Status:     constant.SensorReadingStatus_Pending,
		RetryCount: 0,
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
