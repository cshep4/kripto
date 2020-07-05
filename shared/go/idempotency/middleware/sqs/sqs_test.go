package sqs_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cshep4/kripto/shared/go/idempotency"
	"github.com/cshep4/kripto/shared/go/idempotency/internal/mocks/idempotency"
	"github.com/cshep4/kripto/shared/go/idempotency/middleware/sqs"
	"github.com/cshep4/kripto/shared/go/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testError string

func (t testError) Error() string {
	return string(t)
}

func TestNewMiddleware(t *testing.T) {
	t.Run("returns error if idempotencer is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s, err := sqs.NewMiddleware(nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(sqs.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "idempotencer", ipErr.Parameter)
	})

	t.Run("returns middleware", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		idempotencer := idempotency_mocks.NewMockIdempotencer(ctrl)

		s, err := sqs.NewMiddleware(idempotencer)
		require.NoError(t, err)
		assert.NotNil(t, s)
	})
}

func TestMiddleware_PreExecute(t *testing.T) {
	t.Run("returns error if cannot decode request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		idempotencer := idempotency_mocks.NewMockIdempotencer(ctrl)

		s, err := sqs.NewMiddleware(idempotencer)
		require.NoError(t, err)

		done, resCtx, res, err := s.PreExecute(context.Background(), []byte("❌"))
		require.Error(t, err)

		assert.False(t, done)
		assert.Nil(t, resCtx)
		assert.Nil(t, res)
	})

	t.Run("returns error if error checking idempotency", func(t *testing.T) {
		const (
			messageId           = "messageId"
			testErr   testError = "error"
		)
		var (
			ctrl, ctx = gomock.WithContext(context.Background(), t)

			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			sqsEvent = events.SQSEvent{
				Records: []events.SQSMessage{{
					MessageId: messageId,
				}},
			}
		)
		defer ctrl.Finish()

		ctx = log.WithServiceName(ctx, log.New("debug"), "sqsIdempotencer")

		payload, err := json.Marshal(sqsEvent)
		require.NoError(t, err)

		s, err := sqs.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().Check(ctx, messageId).Return(nil, testErr)

		done, resCtx, res, err := s.PreExecute(ctx, payload)
		require.Error(t, err)

		assert.False(t, done)
		assert.Nil(t, resCtx)
		assert.Nil(t, res)
		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns error if error marking comeplete", func(t *testing.T) {
		const (
			messageId           = "messageId"
			testErr   testError = "error"
		)
		var (
			ctrl, ctx = gomock.WithContext(context.Background(), t)

			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			sqsEvent = events.SQSEvent{
				Records: []events.SQSMessage{{
					MessageId: messageId,
				}},
			}
		)
		defer ctrl.Finish()

		ctx = log.WithServiceName(ctx, log.New("debug"), "sqsIdempotencer")

		payload, err := json.Marshal(sqsEvent)
		require.NoError(t, err)

		s, err := sqs.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().Check(ctx, messageId).Return(&idempotency.Response{}, nil)
		idempotencer.EXPECT().MarkComplete(ctx, messageId, nil).Return(testErr)

		done, resCtx, res, err := s.PreExecute(ctx, payload)
		require.Error(t, err)

		assert.False(t, done)
		assert.Nil(t, resCtx)
		assert.Nil(t, res)
		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("marks complete and returns messages", func(t *testing.T) {
		const messageId = "messageId"
		var (
			ctrl, ctx = gomock.WithContext(context.Background(), t)

			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			sqsEvent = events.SQSEvent{
				Records: []events.SQSMessage{{
					MessageId: messageId,
				}},
			}
		)
		defer ctrl.Finish()

		ctx = log.WithServiceName(ctx, log.New("debug"), "sqsIdempotencer")

		payload, err := json.Marshal(sqsEvent)
		require.NoError(t, err)

		s, err := sqs.NewMiddleware(idempotencer)
		require.NoError(t, err)

		idempotencer.EXPECT().Check(ctx, messageId).Return(&idempotency.Response{}, nil)
		idempotencer.EXPECT().MarkComplete(ctx, messageId, nil).Return(nil)

		done, resCtx, res, err := s.PreExecute(ctx, payload)
		require.NoError(t, err)

		assert.False(t, done)
		assert.Equal(t, ctx, resCtx)

		var resEvents events.SQSEvent
		err = json.Unmarshal(res, &resEvents)
		require.NoError(t, err)

		assert.Equal(t, sqsEvent, resEvents)
	})

	t.Run("strips out messages that have already been processed", func(t *testing.T) {
		const (
			messageId  = "messageId"
			messageId2 = "messageId 2"
		)
		var (
			ctrl, ctx = gomock.WithContext(context.Background(), t)

			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)

			sqsEvent = events.SQSEvent{
				Records: []events.SQSMessage{
					{MessageId: messageId},
					{MessageId: messageId2},
				},
			}
		)
		defer ctrl.Finish()

		ctx = log.WithServiceName(ctx, log.New("debug"), "sqsIdempotencer")

		payload, err := json.Marshal(sqsEvent)
		require.NoError(t, err)

		s, err := sqs.NewMiddleware(idempotencer)
		require.NoError(t, err)

		gomock.InOrder(
			idempotencer.EXPECT().Check(ctx, messageId).Return(&idempotency.Response{Exists: true}, nil),
			idempotencer.EXPECT().Check(ctx, messageId2).Return(&idempotency.Response{}, nil),
			idempotencer.EXPECT().MarkComplete(ctx, messageId2, nil).Return(nil),
		)

		done, resCtx, res, err := s.PreExecute(ctx, payload)
		require.NoError(t, err)

		assert.False(t, done)
		assert.Equal(t, ctx, resCtx)

		var resEvents events.SQSEvent
		err = json.Unmarshal(res, &resEvents)
		require.NoError(t, err)

		assert.Len(t, resEvents.Records, 1)
		assert.Equal(t, sqsEvent.Records[1], resEvents.Records[0])
	})
}
