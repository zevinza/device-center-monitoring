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
)

type Module struct {
	Controller Controller
	Domain     Domain
	Repository Repository
	Service    Service
}

type Controller struct {
	DeviceController       *devicectrl.DeviceController
	SensorController       *sensorctrl.SensorController
	SensorIngestController *sensorctrl.SensorIngestController
}

type Domain struct {
	DeviceDomain        devicedomain.DeviceDomain
	SensorDomain        sensordomain.SensorDomain
	SensorReadingDomain sensorreadingdomain.SensorReadingDomain
}

type Repository struct {
	DeviceRepository        devicerepo.DeviceRepository
	SensorRepository        sensorrepo.SensorRepository
	SensorReadingRepository sensorreadingrepo.SensorReadingRepository
}

type Service struct {
	QueueConsumerService service.QueueConsumer
}
