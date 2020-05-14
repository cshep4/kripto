package lambda_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cshep4/kripto/shared/go/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aws "github.com/aws/aws-lambda-go/lambda"
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/cshep4/kripto/shared/go/lambda/internal/mocks/aws"
	"github.com/golang/mock/gomock"
)

func TestRunner_Invoke(t *testing.T) {
	t.Run("each executor runs in correct order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			first  = aws_mocks.NewMockHandler(ctrl)
			second = aws_mocks.NewMockHandler(ctrl)
			third  = aws_mocks.NewMockHandler(ctrl)
			fourth = aws_mocks.NewMockHandler(ctrl)
			fifth  = aws_mocks.NewMockHandler(ctrl)

			ctx     = context.Background()
			payload = []byte{1}

			function = func() { third.Invoke(ctx, payload) }
		)

		gomock.InOrder(
			first.EXPECT().Invoke(ctx, payload).Return([]byte{}, nil),
			second.EXPECT().Invoke(ctx, payload).Return([]byte{}, nil),
			third.EXPECT().Invoke(ctx, payload).Return([]byte{}, nil),
			fourth.EXPECT().Invoke(ctx, payload).Return([]byte{}, nil),
			fifth.EXPECT().Invoke(ctx, payload).Return([]byte{}, nil),
		)

		runner := lambda.New(
			function,
			lambda.WithPostExecute(postExecutorFunc(fourth)),
			lambda.WithPreExecute(preExecutorFunc(second)),
			lambda.WithPreExecute(preExecutorFunc(first)),
			lambda.WithPostExecute(postExecutorFunc(fifth)),
		)

		_, err := runner.Invoke(ctx, payload)
		require.NoError(t, err)
	})

	t.Run("returns error without executing function if error in pre executor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			first  = aws_mocks.NewMockHandler(ctrl)
			second = aws_mocks.NewMockHandler(ctrl)

			ctx     = context.Background()
			payload = []byte{1}
			testErr = errors.New("error")

			function = func() ([]byte, error) { return second.Invoke(ctx, payload) }
		)

		first.EXPECT().Invoke(ctx, payload).Return(nil, testErr)
		second.EXPECT().Invoke(ctx, payload).Times(0)

		runner := lambda.New(
			function,
			lambda.WithPreExecute(preExecutorFunc(first)),
		)

		_, err := runner.Invoke(ctx, payload)
		require.Error(t, err)

		assert.Equal(t, testErr, err)
	})

	t.Run("post executor not called if error in function", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			first  = aws_mocks.NewMockHandler(ctrl)
			second = aws_mocks.NewMockHandler(ctrl)

			ctx     = context.Background()
			payload = []byte{1}
			testErr = errors.New("error")

			function = func() ([]byte, error) { return first.Invoke(ctx, payload) }
		)

		first.EXPECT().Invoke(ctx, payload).Return(nil, testErr)
		second.EXPECT().Invoke(ctx, payload).Times(0)

		runner := lambda.New(
			function,
			lambda.WithPostExecute(postExecutorFunc(second)),
		)

		_, err := runner.Invoke(ctx, payload)
		require.Error(t, err)

		assert.Equal(t, testErr, err)
	})

	t.Run("returns error if error in post executor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			first  = aws_mocks.NewMockHandler(ctrl)
			second = aws_mocks.NewMockHandler(ctrl)

			ctx     = context.Background()
			payload = []byte{1}
			testErr = errors.New("error")

			function = func() ([]byte, error) { return first.Invoke(ctx, payload) }
		)

		gomock.InOrder(
			first.EXPECT().Invoke(ctx, payload).Return(nil, nil),
			second.EXPECT().Invoke(ctx, payload).Return(nil, testErr),
		)

		runner := lambda.New(
			function,
			lambda.WithPostExecute(postExecutorFunc(second)),
		)

		_, err := runner.Invoke(ctx, payload)
		require.Error(t, err)

		assert.Equal(t, testErr, err)
	})

	t.Run("returns successful response from function", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			first  = aws_mocks.NewMockHandler(ctrl)
			second = aws_mocks.NewMockHandler(ctrl)
			third  = aws_mocks.NewMockHandler(ctrl)

			ctx      = context.Background()
			payload  = []byte{1}
			response = []byte{2}

			function     = func() ([]byte, error) { return second.Invoke(ctx, payload) }
			mockFunction = func() ([]byte, error) { return response, nil }
		)

		gomock.InOrder(
			first.EXPECT().Invoke(ctx, payload).Return(nil, nil),
			second.EXPECT().Invoke(ctx, payload).Return(response, nil),
			third.EXPECT().Invoke(ctx, payload).Return(nil, nil),
		)

		runner := lambda.New(
			function,
			lambda.WithPreExecute(preExecutorFunc(first)),
			lambda.WithPostExecute(postExecutorFunc(third)),
		)

		resp, err := runner.Invoke(ctx, payload)
		require.NoError(t, err)

		expectedResponse, err := aws.NewHandler(mockFunction).Invoke(ctx, payload)
		require.NoError(t, err)

		assert.Equal(t, expectedResponse, resp)
	})

	t.Run("logger is injected into context in pre executor", func(t *testing.T) {
		var (
			ctx = context.Background()

			payload = []byte{1}

			logLevel    = "debug"
			serviceName = "test"
		)

		function := func(ctx context.Context) ([]byte, error) {
			log.Debug(ctx, "🗣")
			return nil, nil
		}

		runner := lambda.New(
			function,
			lambda.WithPreExecute(lambda.LogMiddleware(logLevel, serviceName)),
		)

		_, err := runner.Invoke(ctx, payload)
		require.NoError(t, err)
	})
}

func preExecutorFunc(h aws.Handler) func(ctx context.Context, payload []byte) (context.Context, []byte, error) {
	return func(ctx context.Context, payload []byte) (context.Context, []byte, error) {
		_, err := h.Invoke(ctx, payload)
		return ctx, payload, err
	}
}

func postExecutorFunc(h aws.Handler) func(ctx context.Context, payload []byte) error {
	return func(ctx context.Context, payload []byte) error {
		_, err := h.Invoke(ctx, payload)
		return err
	}
}
