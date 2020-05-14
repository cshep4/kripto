package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(ctx context.Context) (*mongo.Client, error) {
	uri, ok := os.LookupEnv("MONGO_URI")
	if !ok {
		return nil, errors.New("mongo_uri_not_set")
	}

	opts := options.Client().
		SetConnectTimeout(2 * time.Second).
		SetServerSelectionTimeout(2 * time.Second).
		ApplyURI(uri)

	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("opts_validate: %v", err)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}

	return client, nil
}
