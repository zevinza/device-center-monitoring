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
	Device       *devicectrl.DeviceController
	Sensor       *sensorctrl.SensorController
	SensorIngest *sensorctrl.SensorIngestController
}

type Domain struct {
	Device        devicedomain.DeviceDomain
	Sensor        sensordomain.SensorDomain
	SensorReading sensorreadingdomain.SensorReadingDomain
}

type Repository struct {
	Device        devicerepo.DeviceRepository
	Sensor        sensorrepo.SensorRepository
	SensorReading sensorreadingrepo.SensorReadingRepository
}

type Service struct {
	QueueConsumer service.QueueConsumer
}
