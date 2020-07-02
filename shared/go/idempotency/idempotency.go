package idempotency

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	fourDays             = int32(345600)
	collection           = "idempotency"
	maxIdempotencyChecks = 5

	InProgress State = "in_progress"
	Complete   State = "complete"
	Error      State = "error"
)

var ErrMaxAttemptsExceeded = errors.New("max attempts exceeded")

type (
	State string

	idempotencyDoc struct {
		Id        string    `bson:"_id"`
		State     State     `bson:"state"`
		Response  []byte    `bson:"response,omitempty"`
		Err       string    `bson:"error,omitempty"`
		CreatedAt time.Time `bson:"createdAt"`
		UpdatedAt time.Time `bson:"updatedAt"`
	}

	Response struct {
		Exists   bool
		Response []byte `bson:"response,omitempty"`
		Err      error  `bson:"error,omitempty"`
	}

	Idempotencer interface {
		Check(ctx context.Context, key string) (*Response, error)
		MarkComplete(ctx context.Context, key string, response []byte) error
		MarkError(ctx context.Context, key string, err error) error
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

func (m *mongoIdempotencer) Check(ctx context.Context, key string) (*Response, error) {
	doc, err := m.get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	if doc == nil {
		if err := m.store(ctx, key); err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return &Response{
			Exists: false,
		}, nil
	}

	if doc.State == InProgress {
		doc, err = m.waitForResponse(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("wait_for_response: %w", err)
		}
	}

	if doc.State == Error {
		return &Response{
			Exists: true,
			Err:    errors.New(doc.Err),
		}, nil
	}

	return &Response{
		Exists:   true,
		Response: doc.Response,
	}, nil
}

func (m *mongoIdempotencer) waitForResponse(ctx context.Context, key string) (*idempotencyDoc, error) {
	ticker := time.NewTicker(200 * time.Millisecond)
	for i := 0; i < maxIdempotencyChecks; i++ {
		select {
		case <-ctx.Done():
			return nil, errors.New("context cancelled")
		case <-ticker.C:
			doc, err := m.get(ctx, key)
			if err != nil {
				return nil, fmt.Errorf("get: %w", err)
			}
			if doc.State != InProgress {
				return doc, err
			}
		}
	}

	return nil, ErrMaxAttemptsExceeded
}

func (m *mongoIdempotencer) get(ctx context.Context, key string) (*idempotencyDoc, error) {
	var doc idempotencyDoc
	err := m.client.
		Database(m.database).
		Collection(collection).
		FindOne(ctx, bson.D{{Key: "_id", Value: key}}).
		Decode(&doc)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, nil
		default:
			return nil, fmt.Errorf("find_one: %w", err)
		}
	}

	return &doc, nil
}

func (m *mongoIdempotencer) store(ctx context.Context, key string) error {
	_, err := m.client.
		Database(m.database).
		Collection(collection).
		InsertOne(ctx, &idempotencyDoc{
			Id:        key,
			State:     InProgress,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	if err != nil {
		return fmt.Errorf("insert_one: %w", err)
	}

	return nil
}

func (m *mongoIdempotencer) MarkComplete(ctx context.Context, key string, response []byte) error {
	return m.update(ctx, key, bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "response", Value: response},
			{Key: "state", Value: Complete},
			{Key: "updatedAt", Value: time.Now()},
		},
	}})
}

func (m *mongoIdempotencer) MarkError(ctx context.Context, key string, err error) error {
	return m.update(ctx, key, bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "error", Value: err.Error()},
			{Key: "state", Value: Error},
			{Key: "updatedAt", Value: time.Now()},
		},
	}})
}

func (m *mongoIdempotencer) update(ctx context.Context, key string, update bson.D) error {
	res, err := m.client.
		Database(m.database).
		Collection(collection).
		UpdateOne(ctx, bson.D{{Key: "_id", Value: key}}, update)
	if err != nil {
		return fmt.Errorf("update_one: %w", err)
	}

	if res.MatchedCount == 0 {
		return errors.New("item not found")
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
