package initialization

import (
	"api/app/master-service/controller/devicectrl"
	"api/app/master-service/controller/sensorctrl"
	"api/app/master-service/domain/devicedomain"
	"api/app/master-service/domain/sensordomain"
	"api/app/master-service/domain/sensorreadingdomain"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorreadingrepo"
	"api/app/master-service/repository/sensorrepo"
	"api/app/master-service/service"
	"api/services/cache"
	"api/services/database"
	"api/services/queue"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init() Module {
	mongoClient := database.Mongo
	mdb := database.MongoDatabase
	Redis := cache.Redis

	repositories := InitRepositories(mdb)
	domains := InitDomains(repositories, mongoClient, Redis)
	controllers := InitControllers(domains)
	services := InitServices(repositories, Redis)

	return Module{
		Controller: controllers,
		Domain:     domains,
		Repository: repositories,
		Service:    services,
	}
}

func InitControllers(domains Domain) Controller {
	return Controller{
		Device:       devicectrl.New(domains.Device),
		Sensor:       sensorctrl.New(domains.Sensor),
		SensorIngest: sensorctrl.NewSensorIngestController(domains.SensorReading),
	}
}

func InitDomains(repositories Repository, db *mongo.Client, rdb *redis.Client) Domain {
	// Initialize queue for sensor reading domain
	qName := viper.GetString("REDIS_QUEUE_NAME")
	if qName == "" {
		qName = "sensor_data_queue"
	}
	dlqName := viper.GetString("REDIS_DLQ_NAME")
	if dlqName == "" {
		dlqName = "sensor_data_dlq"
	}
	queue := queue.NewRedisQueue(rdb, qName, dlqName)

	return Domain{
		Device:        devicedomain.New(repositories.Device),
		Sensor:        sensordomain.New(repositories.Sensor, repositories.Device),
		SensorReading: sensorreadingdomain.New(repositories.SensorReading, repositories.Device, repositories.Sensor, queue),
	}
}

func InitRepositories(mdb *mongo.Database) Repository {
	return Repository{
		Device:        devicerepo.New(mdb),
		Sensor:        sensorrepo.New(mdb),
		SensorReading: sensorreadingrepo.New(mdb),
	}
}

func InitServices(repositories Repository, rdb *redis.Client) Service {
	return Service{
		QueueConsumer: service.NewQueueConsumer(repositories.SensorReading, rdb),
	}
}
