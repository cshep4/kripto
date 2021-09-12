//+build integration

package mongo_test

import (
	"context"
	"testing"
	"time"

	store "github.com/cshep4/kripto/services/data-storer/internal/store/rate/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				Database("rate").
				Drop(ctx)
			require.NoError(t, err)

			err = client.Disconnect(ctx)
			require.NoError(t, err)
		})

		_, err := client.
			Database("rate").
			Collection("rate").
			Indexes().
			CreateOne(
				ctx,
				mongo.IndexModel{
					Keys:    bsonx.Doc{{Key: "dateTime", Value: bsonx.Int64(1)}},
					Options: options.Index().SetName("dateTimeIdx"),
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
				Database("rate").
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

func TestStore_GetPreviousMonth(t *testing.T) {
	t.Run("get all rates from the past month", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		store, err := store.New(ctx, client)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.
				Database("rate").
				Drop(ctx)
			require.NoError(t, err)

			err = store.Close(ctx)
			require.NoError(t, err)
		})

		var (
			now         = time.Now().Round(time.Second).UTC()
			yesterday   = time.Now().AddDate(0, 0, -1).Round(time.Second).UTC()
			lastMonth   = time.Now().AddDate(0, -1, -8).Round(time.Second).UTC()
			monthBefore = time.Now().AddDate(0, -2, -15).Round(time.Second).UTC()
		)
		const (
			one   = float64(1)
			two   = float64(2)
			three = float64(3)
			four  = float64(4)
		)

		err = store.Store(ctx, one, now)
		require.NoError(t, err)

		err = store.Store(ctx, two, monthBefore)
		require.NoError(t, err)

		err = store.Store(ctx, three, lastMonth)
		require.NoError(t, err)

		err = store.Store(ctx, four, yesterday)
		require.NoError(t, err)

		rates, err := store.GetPreviousMonth(ctx)
		require.NoError(t, err)

		assert.Len(t, rates, 2)

		assert.Equal(t, one, rates[0].Rate)
		assert.Equal(t, now, rates[0].DateTime)

		assert.Equal(t, four, rates[1].Rate)
		assert.Equal(t, yesterday, rates[1].DateTime)
	})
}

func TestStore_Store(t *testing.T) {
	t.Run("stores rate in db", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		store, err := store.New(ctx, client)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.
				Database("rate").
				Drop(ctx)
			require.NoError(t, err)

			err = store.Close(ctx)
			require.NoError(t, err)
		})

		const rate = 1234.124
		now := time.Now().Round(time.Second).UTC()

		err = store.Store(ctx, rate, now)
		require.NoError(t, err)

		rates, err := store.GetPreviousMonth(ctx)
		require.NoError(t, err)

		assert.Len(t, rates, 1)
		assert.Equal(t, rate, rates[0].Rate)
		assert.Equal(t, now, rates[0].DateTime)
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
