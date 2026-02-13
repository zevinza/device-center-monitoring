package initialization

import (
	"api/app/master-service/controller/devicectrl"
	"api/app/master-service/controller/sensorcategoryctrl"
	"api/app/master-service/controller/sensorctrl"
	"api/app/master-service/domain/devicedomain"
	"api/app/master-service/domain/sensorcategorydomain"
	"api/app/master-service/domain/sensordomain"
	"api/app/master-service/domain/sensorreadingdomain"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorcategory"
	"api/app/master-service/repository/sensorreadingrepo"
	"api/app/master-service/repository/sensorrepo"
	"api/app/master-service/service"
	"api/services/cache"
	"api/services/database"
	"api/services/queue"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Init() Module {
	mdb := database.MongoDatabase
	pgdb := database.DB
	Redis := cache.Redis

	repositories := InitRepositories(mdb)
	domains := InitDomains(pgdb, Redis, repositories)
	services := InitServices(Redis, pgdb, repositories)
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
		DeviceController:         devicectrl.New(domains.DeviceDomain),
		SensorController:         sensorctrl.New(domains.SensorDomain),
		SensorCategoryController: sensorcategoryctrl.New(domains.SensorCategoryDomain),
		SensorIngestController:   sensorctrl.NewSensorIngestController(domains.SensorReadingDomain),
	}
}

func InitDomains(pgdb *gorm.DB, rdb *redis.Client, repositories Repository) Domain {
	// Initialize queue for sensor reading domain
	queue := queue.NewRedisQueue(rdb, viper.GetString("REDIS_QUEUE_NAME"), viper.GetString("REDIS_DLQ_NAME"))

	return Domain{
		DeviceDomain:         devicedomain.New(pgdb, repositories.DeviceRepository),
		SensorDomain:         sensordomain.New(pgdb, repositories.SensorRepository, repositories.DeviceRepository, repositories.SensorReadingRepository, repositories.SensorCategoryRepository),
		SensorCategoryDomain: sensorcategorydomain.New(pgdb, repositories.SensorCategoryRepository),
		SensorReadingDomain:  sensorreadingdomain.New(pgdb, queue, repositories.SensorReadingRepository, repositories.DeviceRepository, repositories.SensorRepository),
	}
}

func InitRepositories(mdb *mongo.Database) Repository {
	return Repository{
		DeviceRepository:         devicerepo.New(),
		SensorRepository:         sensorrepo.New(),
		SensorCategoryRepository: sensorcategory.New(),
		SensorReadingRepository:  sensorreadingrepo.New(mdb),
	}
}

func InitServices(rdb *redis.Client, pgdb *gorm.DB, repositories Repository) Service {
	return Service{
		QueueConsumerService: service.NewQueueConsumer(
			pgdb,
			rdb,
			repositories.SensorReadingRepository,
			repositories.DeviceRepository,
			repositories.SensorRepository,
		),
	}
}
