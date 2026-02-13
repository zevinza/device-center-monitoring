package sensorrepo

import (
	"api/app/master-service/model"
	"api/constant"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensorRepository interface {
	GetByDeviceID(ctx context.Context, deviceID primitive.ObjectID, limit int64) ([]model.Sensor, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Sensor, error)
	Create(ctx context.Context, sensor *model.Sensor) error
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*model.Sensor, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type sensorRepository struct {
	coll *mongo.Collection
}

func New(db *mongo.Database) SensorRepository {
	return &sensorRepository{coll: db.Collection(constant.Collection_Sensors)}
}

func (r *sensorRepository) GetByDeviceID(ctx context.Context, deviceID primitive.ObjectID, limit int64) ([]model.Sensor, error) {
	if limit <= 0 {
		limit = 100
	}
	cur, err := r.coll.Find(ctx, bson.M{"device_id": deviceID}, options.Find().SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]model.Sensor, 0)
	for cur.Next(ctx) {
		var s model.Sensor
		if err := cur.Decode(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, cur.Err()
}

func (r *sensorRepository) Create(ctx context.Context, sensor *model.Sensor) error {
	now := time.Now().UTC()
	if sensor.ID.IsZero() {
		sensor.ID = primitive.NewObjectID()
	}
	if sensor.CreatedAt.IsZero() {
		sensor.CreatedAt = now
	}
	sensor.UpdatedAt = now

	_, err := r.coll.InsertOne(ctx, sensor)
	return err
}

func (r *sensorRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*model.Sensor, error) {
	update["updated_at"] = time.Now().UTC()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("sensor not found")
		}
		return nil, err
	}

	var out model.Sensor
	if err := res.Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *sensorRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("sensor not found")
	}
	return nil
}

func (r *sensorRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Sensor, error) {
	var out model.Sensor
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("sensor not found")
		}
		return nil, err
	}
	return &out, nil
}
