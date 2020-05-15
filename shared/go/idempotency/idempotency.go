package idempotency

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

const (
	fourDays   = int32(345600)
	collection = "idempotency"
)

type (
	idempotencyDoc struct {
		Id        string    `bson:"_id"`
		CreatedAt time.Time `bson:"createdAt"`
	}
	Idempotencer interface {
		Check(ctx context.Context, key string) (bool, error)
	}

	mongoIdempotencer struct {
		database string
		client   *mongo.Client
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func New(ctx context.Context, database string, client *mongo.Client) (*mongoIdempotencer, error) {
	switch {
	case client == nil:
		return nil, InvalidParameterError{Parameter: "client"}
	case database == "":
		return nil, InvalidParameterError{Parameter: "database"}
	}

	i := &mongoIdempotencer{
		database: database,
		client:   client,
	}

	if err := i.ping(ctx); err != nil {
		return nil, err
	}

	if err := i.ensureIndexes(ctx); err != nil {
		return nil, err
	}

	return i, nil
}

func (m *mongoIdempotencer) ensureIndexes(ctx context.Context) error {
	_, err := m.client.
		Database(m.database).
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
				SetBackground(true).
				SetExpireAfterSeconds(fourDays),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoIdempotencer) Check(ctx context.Context, key string) (bool, error) {
	exists, err := m.getKey(ctx, key)
	if err != nil {
		return false, fmt.Errorf("check_idempotency: %w", err)
	}
	if exists {
		return true, nil
	}

	if err := m.storeKey(ctx, key); err != nil {
		return false, fmt.Errorf("store_idempotency: %w", err)
	}

	return false, nil
}

func (m *mongoIdempotencer) getKey(ctx context.Context, key string) (bool, error) {
	filter := bson.D{{Key: "_id", Value: key}}

	var doc idempotencyDoc
	err := m.client.
		Database(m.database).
		Collection(collection).
		FindOne(ctx, filter).
		Decode(&doc)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return true, nil
		default:
			return false, fmt.Errorf("find_one: %w", err)
		}
	}

	return false, nil
}

func (m *mongoIdempotencer) storeKey(ctx context.Context, key string) error {
	_, err := m.client.
		Database(m.database).
		Collection(collection).
		InsertOne(ctx, &idempotencyDoc{
			Id:        key,
			CreatedAt: time.Now(),
		})
	if err != nil {
		return fmt.Errorf("insert_one: %w", err)
	}

	return nil
}

func (m *mongoIdempotencer) ping(ctx context.Context) error {
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)
	return m.client.Ping(ctx, nil)
}

func (m *mongoIdempotencer) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
