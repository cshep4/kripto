package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	db         = "trade"
	collection = "trade"
)

type (
	store struct {
		client *mongo.Client
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func New(ctx context.Context, client *mongo.Client) (*store, error) {
	if client == nil {
		return nil, InvalidParameterError{Parameter: "client"}
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
				{Key: "createdAt", Value: bsonx.Int64(1)},
			},
			Options: options.Index().
				SetName("createdAtIdx").
				SetUnique(true).
				SetBackground(true),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) Store(ctx context.Context, trade model.Trade) error {
	t, err := fromTrade(trade)
	if err != nil {
		return fmt.Errorf("map_document: %w", err)
	}

	_, err = s.client.
		Database(db).
		Collection(collection).
		InsertOne(ctx, t)
	if err != nil {
		return fmt.Errorf("insert_one: %w", err)
	}

	return nil
}

func (s *store) GetPreviousWeeks(ctx context.Context) ([]model.Trade, error) {
	panic("implement me")
}

func (s *store) ping(ctx context.Context) error {
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)
	return s.client.Ping(ctx, nil)
}

func (s *store) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
