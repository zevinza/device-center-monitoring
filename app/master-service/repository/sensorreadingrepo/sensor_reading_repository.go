package sensorreadingrepo

import (
	"api/app/master-service/model"
	"api/constant"
	"api/lib"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensorReadingRepository interface {
	Create(ctx context.Context, reading *model.SensorReading) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.SensorReading, error)
	GetBySensorID(ctx context.Context, sensorID string, limit int64) ([]model.SensorReading, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
	IncrementRetryCount(ctx context.Context, id primitive.ObjectID) (int, error)
}

type sensorReadingRepository struct {
	coll *mongo.Collection
}

func New(db *mongo.Database) SensorReadingRepository {
	return &sensorReadingRepository{coll: db.Collection(constant.Collection_SensorReadings)}
}

func (r *sensorReadingRepository) Create(ctx context.Context, reading *model.SensorReading) error {
	if reading.ID.IsZero() {
		reading.ID = primitive.NewObjectID()
	}
	reading.CreatedAt = lib.TimeNow()
	reading.UpdatedAt = lib.TimeNow()

	_, err := r.coll.InsertOne(ctx, reading)
	return err
}

func (r *sensorReadingRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.SensorReading, error) {
	var out model.SensorReading
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("sensor reading not found")
		}
		return nil, err
	}
	return &out, nil
}

func (r *sensorReadingRepository) GetBySensorID(ctx context.Context, sensorID string, limit int64) ([]model.SensorReading, error) {
	var readings []model.SensorReading
	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.M{"timestamp": -1, "created_at": -1})
	cursor, err := r.coll.Find(ctx, bson.M{"sensor_id": sensorID}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("sensor readings not found")
		}
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var reading model.SensorReading
		if err := cursor.Decode(&reading); err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	return readings, nil
}

func (r *sensorReadingRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	res, err := r.coll.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("sensor reading not found")
	}
	return nil
}

func (r *sensorReadingRepository) IncrementRetryCount(ctx context.Context, id primitive.ObjectID) (int, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$inc": bson.M{"retry_count": 1},
		"$set": bson.M{"updated_at": time.Now().UTC()},
	}, opts)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, fmt.Errorf("sensor reading not found")
		}
		return 0, err
	}
	var out model.SensorReading
	if err := res.Decode(&out); err != nil {
		return 0, err
	}
	return out.RetryCount, nil
}
