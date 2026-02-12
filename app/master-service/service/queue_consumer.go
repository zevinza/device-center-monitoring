package service

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorreadingrepo"
	"api/app/master-service/repository/sensorrepo"
	"api/constant"
	"api/services/queue"
	"api/utils/httpreq"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueConsumer interface {
	Start(ctx context.Context)
}

type queueConsumer struct {
	readingRepository sensorreadingrepo.SensorReadingRepository
	deviceRepository  devicerepo.DeviceRepository
	sensorRepository  sensorrepo.SensorRepository
	queue             *queue.RedisQueue
	maxRetries        int
	backoffBase       int
}

func NewQueueConsumer(
	readingRepository sensorreadingrepo.SensorReadingRepository,
	deviceRepository devicerepo.DeviceRepository,
	sensorRepository sensorrepo.SensorRepository,
	rdb *redis.Client,
) QueueConsumer {
	qName := viper.GetString("REDIS_QUEUE_NAME")
	if qName == "" {
		qName = "sensor_data_queue"
	}
	dlqName := viper.GetString("REDIS_DLQ_NAME")
	if dlqName == "" {
		dlqName = "sensor_data_dlq"
	}

	maxRetries := viper.GetInt("MAX_RETRIES")
	if maxRetries <= 0 {
		maxRetries = 3
	}

	backoffBase := viper.GetInt("RETRY_BACKOFF_BASE")
	if backoffBase <= 0 {
		backoffBase = 1
	}

	return &queueConsumer{
		readingRepository: readingRepository,
		queue:             queue.NewRedisQueue(rdb, qName, dlqName),
		maxRetries:        maxRetries,
		backoffBase:       backoffBase,
		deviceRepository:  deviceRepository,
		sensorRepository:  sensorRepository,
	}
}

func (qc *queueConsumer) Start(ctx context.Context) {
	log.Println("Starting queue consumer...")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Queue consumer stopped")
			return
		case <-ticker.C:
			qc.processMessage(ctx)
		}
	}
}

func (qc *queueConsumer) processMessage(ctx context.Context) {
	// Blocking dequeue with 1 second timeout
	data, err := qc.queue.DequeueBlocking(ctx, 1*time.Second)
	if err != nil {
		if err != redis.Nil {
			log.Printf("Error dequeuing message: %v", err)
		}
		return
	}
	if data == nil {
		return
	}

	var msg model.QueueMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Error unmarshaling queue message: %v", err)
		return
	}

	readingID, err := primitive.ObjectIDFromHex(msg.ReadingID)
	if err != nil {
		log.Printf("Invalid reading_id in queue message: %v", err)
		return
	}

	reading, err := qc.readingRepository.GetByID(ctx, readingID)
	if err != nil {
		log.Printf("Error finding reading by id: %v", err)
		return
	}

	sensor, err := qc.sensorRepository.GetByID(ctx, reading.SensorID)
	if err != nil {
		log.Printf("Error finding sensor by id: %v", err)
		return
	}

	device, err := qc.deviceRepository.GetByID(ctx, reading.DeviceID)
	if err != nil {
		log.Printf("Error finding device by id: %v", err)
		return
	}

	result := model.SensorIngestResult{
		SensorReading: *reading,
		Device:        device,
		Sensor:        sensor,
	}

	// Attempt to send to client service
	success := qc.sendToClient(&result)

	if success {
		if err := qc.readingRepository.UpdateStatus(ctx, readingID, constant.SensorReadingStatus_Success); err != nil {
			log.Printf("Error updating reading status to success: %v", err)
		} else {
			log.Printf("Successfully delivered reading %s to client service", msg.ReadingID)
		}
	} else {
		newRetryCount, err := qc.readingRepository.IncrementRetryCount(ctx, readingID)
		if err != nil {
			log.Printf("Error incrementing retry count: %v", err)
			return
		}

		if newRetryCount >= qc.maxRetries {
			if err := qc.readingRepository.UpdateStatus(ctx, readingID, constant.SensorReadingStatus_Failed); err != nil {
				log.Printf("Error updating reading status to failed: %v", err)
			}
			if err := qc.queue.EnqueueDLQ(ctx, data); err != nil {
				log.Printf("Error enqueueing to DLQ: %v", err)
			} else {
				log.Printf("Reading %s moved to dead letter queue after %d retries", msg.ReadingID, newRetryCount)
			}
		} else {
			backoffSeconds := int(math.Pow(float64(qc.backoffBase), float64(newRetryCount)))
			log.Printf("Retrying reading %s (attempt %d/%d) after %d seconds", msg.ReadingID, newRetryCount, qc.maxRetries, backoffSeconds)

			go func() {
				time.Sleep(time.Duration(backoffSeconds) * time.Second)
				if err := qc.queue.Enqueue(ctx, data); err != nil {
					log.Printf("Error re-enqueueing message: %v", err)
				}
			}()
		}
	}
}

func (qc *queueConsumer) sendToClient(body *model.SensorIngestResult) bool {
	url := fmt.Sprintf("http://%s:%d", viper.GetString("CLIENT_HOST"), viper.GetInt("CLIENT_PORT")) + "/receive"
	resp, err := httpreq.NewClient().
		Url(url).
		Headers(map[string]string{
			"Content-Type": "application/json",
			"X-API-Key":    viper.GetString("SERVER_SECRET_KEY"),
		}).Body(body).Post()
	if err != nil {
		log.Printf("Error sending to client service: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	}

	log.Printf("Client service returned non-2xx status: %d", resp.StatusCode)
	return false
}
