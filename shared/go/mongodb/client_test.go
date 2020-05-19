//+build integration

package mongodb_test

import (
	"os"
	"testing"

	"github.com/cshep4/kripto/shared/go/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestNew(t *testing.T) {
	t.Run("returns error if MONGO_URI env variable is not set", func(t *testing.T) {
		err := os.Unsetenv("MONGO_URI")
		require.NoError(t, err)
		
		client, err := mongodb.New(context.Background())
		require.Error(t, err)

		assert.Nil(t, client)
		assert.Equal(t, "mongo_uri_not_set", err.Error())
	})

	t.Run("returns error if MONGO_URI env variable is not valid", func(t *testing.T) {
		err := os.Setenv("MONGO_URI", "invalid")
		require.NoError(t, err)

		t.Cleanup(func() {
			err := os.Unsetenv("MONGO_URI")
			require.NoError(t, err)
		})

		client, err := mongodb.New(context.Background())
		require.Error(t, err)

		assert.Nil(t, client)
	})

	t.Run("returns error if cannot connect to mongo instance", func(t *testing.T) {
		err := os.Setenv("MONGO_URI", "mongod://not-real-host:27017")
		require.NoError(t, err)

		t.Cleanup(func() {
			err := os.Unsetenv("MONGO_URI")
			require.NoError(t, err)
		})

		client, err := mongodb.New(context.Background())
		require.Error(t, err)

		assert.Nil(t, client)
	})

	t.Run("returns mongo client", func(t *testing.T) {
		err := os.Setenv("MONGO_URI", "mongodb://localhost:27017")
		require.NoError(t, err)

		t.Cleanup(func() {
			err := os.Unsetenv("MONGO_URI")
			require.NoError(t, err)
		})

		ctx := context.Background()

		client, err := mongodb.New(ctx)
		require.NoError(t, err)

		assert.NotNil(t, client)

		err = client.Ping(ctx, nil)
		require.NoError(t, err)

		err = client.Disconnect(ctx)
		require.NoError(t, err)
	})
}
