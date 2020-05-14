package mongo

import (
	"context"
	"errors"
	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	db         = "user"
	collection = "user"
)

type store struct {
	client *mongo.Client
}

func New(ctx context.Context, client *mongo.Client) (*store, error) {
	if client == nil {
		return nil, errors.New("mongo_client_is_nil")
	}

	s := &store{
		client: client,
	}

	if err := s.ping(ctx); err != nil {
		return nil, err
	}

	if err := s.ensureIndexes(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *store) ensureIndexes(ctx context.Context) error {
	_, err := s.client.
		Database(db).
		Collection(collection).
		Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "email", Value: bsonx.Int64(1)},
			},
			Options: options.Index().
				SetName("email_idx").
				SetUnique(true).
				SetBackground(true),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *store) Store(ctx context.Context, user model.User) error {
	panic("implement me")
}

func (s *store) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	panic("implement me")
}

func (s *store) ping(ctx context.Context) error {
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)
	return s.client.Ping(ctx, nil)
}

func (s *store) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
