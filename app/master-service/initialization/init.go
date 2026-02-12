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
	services := InitServices(repositories, Redis)
	controllers := InitControllers(domains)

	return Module{
		Controller: controllers,
		Domain:     domains,
		Repository: repositories,
		Service:    services,
	}
}

func InitControllers(domains Domain) Controller {
	return Controller{
		DeviceController:       devicectrl.New(domains.DeviceDomain),
		SensorController:       sensorctrl.New(domains.SensorDomain),
		SensorIngestController: sensorctrl.NewSensorIngestController(domains.SensorReadingDomain),
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
		DeviceDomain:        devicedomain.New(repositories.DeviceRepository),
		SensorDomain:        sensordomain.New(repositories.SensorRepository, repositories.DeviceRepository),
		SensorReadingDomain: sensorreadingdomain.New(repositories.SensorReadingRepository, repositories.DeviceRepository, repositories.SensorRepository, queue),
	}
}

func InitRepositories(mdb *mongo.Database) Repository {
	return Repository{
		DeviceRepository:        devicerepo.New(mdb),
		SensorRepository:        sensorrepo.New(mdb),
		SensorReadingRepository: sensorreadingrepo.New(mdb),
	}
}

func InitServices(repositories Repository, rdb *redis.Client) Service {
	return Service{
		QueueConsumerService: service.NewQueueConsumer(
			repositories.SensorReadingRepository,
			repositories.DeviceRepository,
			repositories.SensorRepository,
			rdb,
		),
	}
}
