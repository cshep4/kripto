//+build integration

package mongo_test

import (
	"context"
	"testing"
	"time"

	store "github.com/cshep4/kripto/services/data-storer/internal/store/rate/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestStore_GetPreviousWeeks(t *testing.T) {
	t.Run("get all rates from the past week", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		store, err := store.New(ctx, client)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := client.
				Database("rate").
				Collection("rate").
				DeleteMany(ctx, bson.M{})
			require.NoError(t, err)

			err = store.Close(ctx)
			require.NoError(t, err)
		})

		var (
			now        = time.Now().Round(time.Second).UTC()
			yesterday  = time.Now().AddDate(0, 0, -1).Round(time.Second).UTC()
			lastWeek   = time.Now().AddDate(0, 0, -8).Round(time.Second).UTC()
			weekBefore = time.Now().AddDate(0, 0, -15).Round(time.Second).UTC()
		)
		const (
			one   = float64(1)
			two   = float64(2)
			three = float64(3)
			four  = float64(4)
		)

		err = store.Store(ctx, one, now)
		require.NoError(t, err)

		err = store.Store(ctx, two, weekBefore)
		require.NoError(t, err)

		err = store.Store(ctx, three, lastWeek)
		require.NoError(t, err)

		err = store.Store(ctx, four, yesterday)
		require.NoError(t, err)

		rates, err := store.GetPreviousWeeks(ctx)
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
			_, err := client.
				Database("rate").
				Collection("rate").
				DeleteMany(ctx, bson.M{})
			require.NoError(t, err)

			err = store.Close(ctx)
			require.NoError(t, err)
		})

		const rate = 1234.124
		now := time.Now().Round(time.Second).UTC()

		err = store.Store(ctx, rate, now)
		require.NoError(t, err)

		rates, err := store.GetPreviousWeeks(ctx)
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
