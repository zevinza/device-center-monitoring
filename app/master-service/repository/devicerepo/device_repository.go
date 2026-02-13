package devicerepo

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

type DeviceRepository interface {
	FindAll(ctx context.Context, limit int64) ([]model.Device, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Device, error)
	FindByCode(ctx context.Context, deviceCode string) (*model.Device, error)
	Create(ctx context.Context, device *model.Device) error
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*model.Device, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type deviceRepository struct {
	coll *mongo.Collection
}

func New(db *mongo.Database) DeviceRepository {
	return &deviceRepository{coll: db.Collection(constant.Collection_Devices)}
}

func (r *deviceRepository) Create(ctx context.Context, device *model.Device) error {
	now := time.Now().UTC()
	if device.CreatedAt.IsZero() {
		device.CreatedAt = now
	}
	device.UpdatedAt = now

	_, err := r.coll.InsertOne(ctx, device)
	return err
}

func (r *deviceRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*model.Device, error) {
	update["updated_at"] = time.Now().UTC()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("device not found")
		}
		return nil, err
	}

	var out model.Device
	if err := res.Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *deviceRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("device not found")
	}
	return nil
}

func (r *deviceRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Device, error) {
	var out model.Device
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("device not found")
		}
		return nil, err
	}
	return &out, nil
}

func (r *deviceRepository) FindByCode(ctx context.Context, deviceCode string) (*model.Device, error) {
	var out model.Device
	err := r.coll.FindOne(ctx, bson.M{"device_code": deviceCode}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("device not found")
		}
		return nil, err
	}
	return &out, nil
}

func (r *deviceRepository) FindAll(ctx context.Context, limit int64) ([]model.Device, error) {
	if limit <= 0 {
		limit = 100
	}

	cur, err := r.coll.Find(ctx, bson.M{}, options.Find().SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]model.Device, 0)
	for cur.Next(ctx) {
		var d model.Device
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, cur.Err()
}
