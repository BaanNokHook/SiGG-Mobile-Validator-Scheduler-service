package repo

import (
	"context"
	"nextclan/validator-register/mobile-validator-scheduler-service/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeviceRepository interface {
	InsertOne(ctx context.Context, model *entity.Device) (interface{}, error)
	FindById(ctx context.Context, id primitive.ObjectID) (*entity.Device, error)
	FindDeviceByLatestSyncGte(ctx context.Context, time int64, limit int64) ([]*entity.Device, error)
	FindDeviceByLatestSyncLte(ctx context.Context, time int64, limit int64) ([]*entity.Device, error)
	UpdateLatestSyncByDeviceId(ctx context.Context, timestamp int64, model *entity.Device) (bool, error)
	DeleteById(ctx context.Context, id primitive.ObjectID) (int, error)
}

type DeviceRepositoryMongo struct {
	collection *mongo.Collection
}

func NewDeviceRepository(collection *mongo.Collection) *DeviceRepositoryMongo {
	return &DeviceRepositoryMongo{collection}
}

func (r *DeviceRepositoryMongo) InsertOne(ctx context.Context, model *entity.Device) (interface{}, error) {
	result, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}
	model.ID = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID, nil
}

func (r *DeviceRepositoryMongo) FindById(ctx context.Context, id primitive.ObjectID) (*entity.Device, error) {
	var entity entity.Device
	if err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}, options.FindOne().SetSort(bson.M{})).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *DeviceRepositoryMongo) FindDeviceByLatestSyncGte(ctx context.Context, time int64, limit int64) ([]*entity.Device, error) {
	opts := options.Find().SetLimit(limit)
	cur, err := r.collection.Find(context.Background(), bson.M{"latestSync": bson.M{"$gte": time}}, opts)
	if err != nil {
		panic(err)
	}
	defer cur.Close(context.Background())

	results := []*entity.Device{}
	if err = cur.All(context.Background(), &results); err != nil {
		panic(err)
	}
	return results, err
}

func (r *DeviceRepositoryMongo) FindDeviceByLatestSyncLte(ctx context.Context, time int64, limit int64) ([]*entity.Device, error) {
	opts := options.Find().SetLimit(limit)
	cur, err := r.collection.Find(context.Background(), bson.M{"latestSync": bson.M{"$lte": time}}, opts)
	if err != nil {
		panic(err)
	}
	defer cur.Close(context.Background())

	results := []*entity.Device{}
	if err = cur.All(context.Background(), &results); err != nil {
		panic(err)
	}
	return results, err
}

func (r *DeviceRepositoryMongo) UpdateLatestSyncByDeviceId(ctx context.Context, timestamp int64, model *entity.Device) (bool, error) {
	result, err := r.collection.UpdateOne(ctx, bson.M{
		"deviceId": model.DeviceId,
	}, bson.M{
		"$set": bson.M{
			"latestSync": timestamp,
		},
	})
	if err != nil {
		return false, err
	}
	model.LatestSync = timestamp
	return result.MatchedCount > 0, err
}

func (r *DeviceRepositoryMongo) DeleteById(ctx context.Context, id primitive.ObjectID) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
