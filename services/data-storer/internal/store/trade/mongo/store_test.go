//+build integration

package mongo_test

import (
	"context"
	"testing"

	"github.com/cshep4/kripto/services/data-storer/internal/model"
	store "github.com/cshep4/kripto/services/data-storer/internal/store/trade/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func TestNew(t *testing.T) {
	t.Run("returns error if mongo client is nil", func(t *testing.T) {
		s, err := store.New(context.Background(), nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(store.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "client", ipErr.Parameter)
	})

	t.Run("returns error if ping fails", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)

		err := client.Disconnect(ctx)
		require.NoError(t, err)

		s, err := store.New(ctx, client)
		require.Error(t, err)

		assert.Nil(t, s)
	})

	t.Run("returns error if error ensuring indexes", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)

		t.Cleanup(func() {
			err := client.
				Database("trade").
				Drop(ctx)
			require.NoError(t, err)

			err = client.Disconnect(ctx)
			require.NoError(t, err)
		})

		_, err := client.
			Database("trade").
			Collection("trade").
			Indexes().
			CreateOne(
				ctx,
				mongo.IndexModel{
					Keys:    bsonx.Doc{{Key: "createdAt", Value: bsonx.Int64(1)}},
					Options: options.Index().SetName("createdAtIdx"),
				},
			)
		require.NoError(t, err)

		s, err := store.New(ctx, client)
		require.Error(t, err)

		assert.Nil(t, s)
	})

	t.Run("returns store", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)

		t.Cleanup(func() {
			err := client.
				Database("trade").
				Drop(ctx)
			require.NoError(t, err)

			err = client.Disconnect(ctx)
			require.NoError(t, err)
		})

		s, err := store.New(ctx, client)
		require.NoError(t, err)

		assert.NotNil(t, s)
	})
}

func TestStore_Store(t *testing.T) {
	t.Run("returns error if trade id is empty", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		store, err := store.New(ctx, client)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.
				Database("trade").
				Drop(ctx)
			require.NoError(t, err)

			err = store.Close(ctx)
			require.NoError(t, err)
		})

		err = store.Store(ctx, model.Trade{})
		require.Error(t, err)

		assert.Contains(t, err.Error(), "invalid_trade_id")
	})

	t.Run("stores trade in db", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		store, err := store.New(ctx, client)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.
				Database("trade").
				Drop(ctx)
			require.NoError(t, err)

			err = store.Close(ctx)
			require.NoError(t, err)
		})

		const tradeId = "🤝"
		trade := model.Trade{
			Id: tradeId,
		}

		err = store.Store(ctx, trade)
		require.NoError(t, err)

		var res map[string]interface{}
		err = client.
			Database("trade").
			Collection("trade").
			FindOne(
				ctx,
				bson.M{
					"_id": tradeId,
				},
			).Decode(&res)
		require.NoError(t, err)

		assert.Equal(t, tradeId, res["_id"])
	})
}

func newClient(t *testing.T, ctx context.Context) *mongo.Client {
	t.Helper()

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	err = client.Connect(ctx)
	require.NoError(t, err)

	return client
}
