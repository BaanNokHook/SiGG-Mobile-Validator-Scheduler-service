package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
}

func New(connectionString string) *MongoDB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		panic(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}
	return &MongoDB{
		Client: client,
	}
}

func (m *MongoDB) Close() {
	err := m.Client.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}
}
