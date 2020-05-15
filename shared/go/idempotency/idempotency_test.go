//+build integration

package idempotency_test

import (
	"context"
	"testing"

	"github.com/cshep4/kripto/shared/go/idempotency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoIdempotencer_Check(t *testing.T) {
	t.Run("returns false if idempotency key not used", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		idempotencer, err := idempotency.New(ctx, "database", client)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := client.
				Database("database").
				Collection("idempotency").
				DeleteMany(ctx, bson.M{})
			require.NoError(t, err)

			err = idempotencer.Close(ctx)
			require.NoError(t, err)
		})

		ok, err := idempotencer.Check(ctx, "key")
		require.NoError(t, err)
		assert.False(t, ok)

		ok, err = idempotencer.Check(ctx, "key2")
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("returns true if idempotency key already used", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		idempotencer, err := idempotency.New(ctx, "database", client)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := client.
				Database("database").
				Collection("idempotency").
				DeleteMany(ctx, bson.M{})
			require.NoError(t, err)

			err = idempotencer.Close(ctx)
			require.NoError(t, err)
		})

		const key = "key"

		ok, err := idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.False(t, ok)

		ok, err = idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.True(t, ok)
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
