package migrations

import "api/app/master-service/model"

// ModelMigrations models to automigrate
var ModelMigrations = []any{
	&model.Device{},
	&model.Sensor{},
	&model.SensorCategory{},
}
