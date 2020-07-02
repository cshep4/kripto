//+build integration

package idempotency_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

		res, err := idempotencer.Check(ctx, "key")
		require.NoError(t, err)
		assert.False(t, res.Exists)

		res, err = idempotencer.Check(ctx, "key2")
		require.NoError(t, err)
		assert.False(t, res.Exists)
	})

	t.Run("returns response if idempotency key already used", func(t *testing.T) {
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
		response := []byte{1, 2, 3}

		res, err := idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.False(t, res.Exists)

		err = idempotencer.MarkComplete(ctx, key, response)
		require.NoError(t, err)

		res, err = idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.True(t, res.Exists)
		assert.Nil(t, res.Err)
		assert.Equal(t, response, res.Response)
	})

	t.Run("returns error if error occurred", func(t *testing.T) {
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
		testErr := errors.New("error")

		res, err := idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.False(t, res.Exists)

		err = idempotencer.MarkError(ctx, key, testErr)
		require.NoError(t, err)

		res, err = idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.True(t, res.Exists)
		assert.Nil(t, res.Response)
		assert.Equal(t, testErr, res.Err)
	})

	t.Run("waits for response to be stored then returns", func(t *testing.T) {
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
		response := []byte{1, 2, 3}

		res, err := idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.False(t, res.Exists)

		go func() {
			time.Sleep(500 * time.Millisecond)
			err = idempotencer.MarkComplete(ctx, key, response)
			require.NoError(t, err)
		}()

		res, err = idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.True(t, res.Exists)
		assert.Nil(t, res.Err)
		assert.Equal(t, response, res.Response)
	})

	t.Run("returns error if max attempts exceeded whilst in progress state", func(t *testing.T) {
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

		res, err := idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.False(t, res.Exists)

		res, err = idempotencer.Check(ctx, key)
		require.Error(t, err)

		assert.Equal(t, idempotency.ErrMaxAttemptsExceeded, err)
	})
}

func TestMongoIdempotencer_MarkComplete(t *testing.T) {
	t.Run("returns error if item does not exist", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		idempotencer, err := idempotency.New(ctx, "database", client)
		require.NoError(t, err)

		err = idempotencer.MarkComplete(ctx, "key", []byte{})
		require.Error(t, err)
	})

	t.Run("marks item as complete", func(t *testing.T) {
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
		response := []byte{1, 2, 3}

		res, err := idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.False(t, res.Exists)

		err = idempotencer.MarkComplete(ctx, key, response)
		require.NoError(t, err)

		res, err = idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.True(t, res.Exists)
		assert.Nil(t, res.Err)
		assert.Equal(t, response, res.Response)
	})
}

func TestMongoIdempotencer_MarkError(t *testing.T) {
	t.Run("returns error if item does not exist", func(t *testing.T) {
		ctx := context.Background()

		client := newClient(t, ctx)
		idempotencer, err := idempotency.New(ctx, "database", client)
		require.NoError(t, err)

		err = idempotencer.MarkError(ctx, "key", errors.New("error"))
		require.Error(t, err)
	})

	t.Run("marks item as error", func(t *testing.T) {
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
		testErr := errors.New("error")

		res, err := idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.False(t, res.Exists)

		err = idempotencer.MarkError(ctx, "key", testErr)
		require.NoError(t, err)

		res, err = idempotencer.Check(ctx, key)
		require.NoError(t, err)
		assert.True(t, res.Exists)
		assert.Nil(t, res.Response)
		assert.Equal(t, testErr, res.Err)
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
