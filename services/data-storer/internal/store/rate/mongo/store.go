package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	db         = "rate"
	collection = "rate"
)

type (
	store struct {
		client     *mongo.Client
		collection *mongo.Collection
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
		client:     client,
		collection: client.Database(db).Collection(collection),
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
	_, err := s.collection.Indexes().
		CreateOne(
			ctx,
			mongo.IndexModel{
				Keys: bsonx.Doc{
					{Key: "dateTime", Value: bsonx.Int64(1)},
				},
				Options: options.Index().
					SetName("dateTimeIdx").
					SetUnique(true).
					SetBackground(true),
			},
		)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) Store(ctx context.Context, r float64, dateTime time.Time) error {
	_, err := s.collection.InsertOne(ctx, &rate{
		Id:       primitive.NewObjectID(),
		Rate:     r,
		DateTime: dateTime,
	})
	if err != nil {
		return fmt.Errorf("insert_one: %w", err)
	}

	return nil
}

func (s *store) GetPreviousWeeks(ctx context.Context) ([]model.Rate, error) {
	cur, err := s.collection.Find(
		ctx,
		bson.D{
			{
				Key: "dateTime",
				Value: bson.D{
					{
						Key:   "$gte",
						Value: time.Now().AddDate(0, 0, -7),
					},
				},
			},
		},
		&options.FindOptions{
			Sort: bson.D{
				bson.E{Key: "dateTime", Value: -1},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	var rates []model.Rate
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var r rate
		err := cur.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}

		rates = append(rates, toRate(r))
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor_err: %w", err)
	}

	return rates, nil
}

func (s *store) ping(ctx context.Context) error {
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)
	return s.client.Ping(ctx, nil)
}

func (s *store) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
