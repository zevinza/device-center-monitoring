package migrations

import "api/app/master-service/model"

var (
	sensorCategory model.SensorCategory
)

func DataSeeds() []any {
	return []any{
		sensorCategory.Seed(),
	}
}
