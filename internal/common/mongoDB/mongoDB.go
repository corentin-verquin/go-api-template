package mongodb

import (
	"context"

	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	*mongo.Client
	Registry *bsoncodec.Registry
	Logger   *zap.Logger
}

func NewMongoDB(cfg *config.Registry, logger *zap.Logger) (*MongoDB, error) {
	registry := bson.NewRegistry()

	clientOptions := options.Client().
		ApplyURI(cfg.String(mongoDBUriKey)).
		SetRegistry(registry)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	return &MongoDB{Client: client, Registry: registry, Logger: logger}, nil
}

func (m *MongoDB) onStart(ctx context.Context) error {
	m.Logger.Info("MongoDB starting...")
	if err := m.Client.Ping(ctx, nil); err != nil {
		m.Logger.Error("Error pinging MongoDB: %v", zap.Error(err))
		return err
	}
	m.Logger.Info("MongoDB started successfully.")
	return nil
}

func (m *MongoDB) onStop(ctx context.Context) error {
	m.Logger.Info("MongoDB stopping...")

	if err := m.Client.Disconnect(ctx); err != nil {
		m.Logger.Error("Error disconnecting MongoDB: %v", zap.Error(err))
		return err
	} else {
		m.Logger.Info("MongoDB disconnected successfully.")
	}
	return nil
}
