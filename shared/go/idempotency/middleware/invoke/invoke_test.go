package invoke_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cshep4/kripto/shared/go/idempotency"
	"github.com/cshep4/kripto/shared/go/idempotency/internal/mocks/idempotency"
	"github.com/cshep4/kripto/shared/go/idempotency/middleware/invoke"
	"github.com/cshep4/kripto/shared/go/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type testError string

func (t testError) Error() string {
	return string(t)
}

func TestNewMiddleware(t *testing.T) {
	t.Run("returns error if idempotencer is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s, err := invoke.NewMiddleware(nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(invoke.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "idempotencer", ipErr.Parameter)
	})

	t.Run("returns middleware", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		idempotencer := idempotency_mocks.NewMockIdempotencer(ctrl)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)
		assert.NotNil(t, s)
	})
}

func TestMiddleware_PreExecute(t *testing.T) {
	t.Run("returns error if cannot decode request", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		idempotencer := idempotency_mocks.NewMockIdempotencer(ctrl)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		done, resCtx, res, err := s.PreExecute(ctx, []byte("❌"))
		require.Error(t, err)

		assert.True(t, done)
		assert.Equal(t, ctx, resCtx)
		assert.Nil(t, res)
	})

	t.Run("returns error if idemepotency key does not exist", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req    = map[string]interface{}{}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		done, resCtx, res, err := s.PreExecute(ctx, b)
		require.Error(t, err)

		assert.True(t, done)
		assert.Equal(t, ctx, resCtx)
		assert.Nil(t, res)
	})

	t.Run("returns error if idemepotency key not valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": []int{1, 2, 3, 4},
			}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		done, resCtx, res, err := s.PreExecute(ctx, b)
		require.Error(t, err)

		assert.True(t, done)
		assert.Equal(t, ctx, resCtx)
		assert.Nil(t, res)
	})

	t.Run("returns error if error checking idempotency", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const (
			idempotencyKey           = "🔑"
			testErr        testError = "error"
		)
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": idempotencyKey,
			}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		ctx = log.WithServiceName(ctx, log.New("debug"), "invokeIdempotencer")

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().Check(ctx, idempotencyKey).Return(nil, testErr)

		done, resCtx, res, err := s.PreExecute(ctx, b)
		require.Error(t, err)

		assert.True(t, done)
		assert.Equal(t, ctx, resCtx)
		assert.Nil(t, res)
		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns request payload if item does not exist", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const (
			idempotencyKey = "🔑"
		)
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": idempotencyKey,
			}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		ctx = log.WithServiceName(ctx, log.New("debug"), "invokeIdempotencer")

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().Check(ctx, idempotencyKey).Return(&idempotency.Response{}, nil)

		done, resCtx, res, err := s.PreExecute(ctx, b)
		require.NoError(t, err)

		assert.False(t, done)
		assert.Equal(t, ctx, resCtx)
		assert.Equal(t, b, res)
	})

	t.Run("returns error if item already exists with error response", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const (
			idempotencyKey           = "🔑"
			testErr        testError = "error"
		)
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": idempotencyKey,
			}
			b, err = json.Marshal(req)

			ir = &idempotency.Response{
				Exists: true,
				Err:    testErr,
			}
		)
		require.NoError(t, err)

		ctx = log.WithServiceName(ctx, log.New("debug"), "invokeIdempotencer")

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().Check(ctx, idempotencyKey).Return(ir, nil)

		done, resCtx, res, err := s.PreExecute(ctx, b)
		require.Error(t, err)

		assert.True(t, done)
		assert.Equal(t, ctx, resCtx)
		assert.Nil(t, res)
		assert.Equal(t, testErr, err)
	})

	t.Run("returns payload if item already exists", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const (
			idempotencyKey = "🔑"
		)
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": idempotencyKey,
			}
			b, err = json.Marshal(req)

			payload = []byte{12, 23}
			ir      = &idempotency.Response{
				Exists:   true,
				Response: payload,
			}
		)
		require.NoError(t, err)

		ctx = log.WithServiceName(ctx, log.New("debug"), "invokeIdempotencer")

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().Check(ctx, idempotencyKey).Return(ir, nil)

		done, resCtx, res, err := s.PreExecute(ctx, b)
		require.NoError(t, err)

		assert.True(t, done)
		assert.Equal(t, ctx, resCtx)
		assert.Equal(t, payload, res)
	})
}

func TestMiddleware_PostExecute(t *testing.T) {
	t.Run("returns no error if cannot decode request", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		idempotencer := idempotency_mocks.NewMockIdempotencer(ctrl)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		err = s.PostExecute(ctx, []byte("❌"), []byte("❌"))
		require.NoError(t, err)
	})

	t.Run("returns no error if idempotency key does not exist", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req    = map[string]interface{}{}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		err = s.PostExecute(ctx, b, b)
		require.NoError(t, err)
	})

	t.Run("returns no error if idemepotency key not valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": []int{1, 2, 3, 4},
			}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		err = s.PostExecute(ctx, b, b)
		require.NoError(t, err)
	})

	t.Run("returns no error if error marking complete", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const (
			idempotencyKey           = "🔑"
			testErr        testError = "error"
		)
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": idempotencyKey,
			}
			b, err = json.Marshal(req)
			res    = []byte{123}
		)
		require.NoError(t, err)

		ctx = log.WithServiceName(ctx, log.New("debug"), "invokeIdempotencer")

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().MarkComplete(ctx, idempotencyKey, res).Return(testErr)

		err = s.PostExecute(ctx, b, res)
		require.NoError(t, err)
	})

	t.Run("returns nil if marked complete successful", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const (
			idempotencyKey = "🔑"
		)
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": idempotencyKey,
			}
			b, err = json.Marshal(req)
			res    = []byte{123}
		)
		require.NoError(t, err)

		ctx = log.WithServiceName(ctx, log.New("debug"), "invokeIdempotencer")

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().MarkComplete(ctx, idempotencyKey, res).Return(nil)

		err = s.PostExecute(ctx, b, res)
		require.NoError(t, err)
	})
}

func TestMiddleware_HandleError(t *testing.T) {
	t.Run("return if cannot decode request", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const testErr testError = "error"
		idempotencer := idempotency_mocks.NewMockIdempotencer(ctrl)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().MarkError(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		s.HandleError(ctx, []byte("❌"), testErr)
	})

	t.Run("return if idempotency key does not exist", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const testErr testError = "error"
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req    = map[string]interface{}{}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().MarkError(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		s.HandleError(ctx, b, testErr)
	})

	t.Run("returns if idemepotency key not valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const testErr testError = "error"
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": []int{1, 2, 3, 4},
			}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().MarkError(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		s.HandleError(ctx, b, testErr)
	})

	t.Run("returns no error if error marking complete", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const (
			idempotencyKey           = "🔑"
			testErr        testError = "error"
		)
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": idempotencyKey,
			}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		ctx = log.WithServiceName(ctx, log.New("debug"), "invokeIdempotencer")

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().MarkError(ctx, idempotencyKey, testErr).Return(testErr)

		s.HandleError(ctx, b, testErr)
	})

	t.Run("returns nil if marked complete successful", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		const (
			idempotencyKey           = "🔑"
			testErr        testError = "error"
		)
		var (
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			req = map[string]interface{}{
				"idempotencyKey": idempotencyKey,
			}
			b, err = json.Marshal(req)
		)
		require.NoError(t, err)

		ctx = log.WithServiceName(ctx, log.New("debug"), "invokeIdempotencer")

		s, err := invoke.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().MarkError(ctx, idempotencyKey, testErr).Return(nil)

		s.HandleError(ctx, b, testErr)
	})
}
