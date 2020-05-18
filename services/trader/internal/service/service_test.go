package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/cshep4/kripto/services/trader/internal/mocks/aws"
	"github.com/cshep4/kripto/services/trader/internal/mocks/trader"
	"github.com/cshep4/kripto/services/trader/internal/service"
	trade "github.com/cshep4/kripto/services/trader/internal/trader"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("returns error if amount is empty", func(t *testing.T) {
		s, err := service.New("", "", nil, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "amount", ipErr.Parameter)
	})

	t.Run("returns error if queueUrl is empty", func(t *testing.T) {
		s, err := service.New("amount", "", nil, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "queueUrl", ipErr.Parameter)
	})

	t.Run("returns error if sqsClient is empty", func(t *testing.T) {
		s, err := service.New("amount", "url", nil, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "sqsClient", ipErr.Parameter)
	})

	t.Run("returns error if trader is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqs := aws_mocks.NewMockSQS(ctrl)

		s, err := service.New("amount", "url", sqs, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "trader", ipErr.Parameter)
	})

	t.Run("returns service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqs := aws_mocks.NewMockSQS(ctrl)
		trader := trader_mocks.NewMockTrader(ctrl)

		s, err := service.New("amount", "url", sqs, trader)
		require.NoError(t, err)

		assert.NotNil(t, s)
	})
}

func TestService_Trade(t *testing.T) {
	t.Run("returns error if error trading", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqs := aws_mocks.NewMockSQS(ctrl)
		trader := trader_mocks.NewMockTrader(ctrl)

		const (
			amount    = "amount"
			url       = "url"
			tradeType = "type"
		)
		testErr := errors.New("error")

		trader.EXPECT().Trade(trade.TradeType(tradeType), amount).Return(nil, testErr)

		s, err := service.New(amount, url, sqs, trader)
		require.NoError(t, err)

		err = s.Trade(context.Background(), tradeType)
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns error if error sending message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqsClient := aws_mocks.NewMockSQS(ctrl)
		trader := trader_mocks.NewMockTrader(ctrl)

		const (
			amount    = "amount"
			url       = "url"
			tradeType = "type"
			id        = "tradeId"
		)
		testErr := errors.New("error")

		order := &trade.TradeResponse{
			ID: id,
		}
		b, err := json.Marshal(order)
		require.NoError(t, err)

		msg := &sqs.SendMessageInput{
			MessageBody: aws.String(string(b)),
			QueueUrl:    aws.String(url),
		}

		trader.EXPECT().Trade(trade.TradeType(tradeType), amount).Return(order, nil)
		sqsClient.EXPECT().SendMessage(msg).Return(nil, testErr)

		s, err := service.New(amount, url, sqsClient, trader)
		require.NoError(t, err)

		err = s.Trade(context.Background(), tradeType)
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns nil if trade executed successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqsClient := aws_mocks.NewMockSQS(ctrl)
		trader := trader_mocks.NewMockTrader(ctrl)

		const (
			amount    = "amount"
			url       = "url"
			tradeType = "type"
			id        = "tradeId"
		)

		order := &trade.TradeResponse{
			ID: id,
		}
		b, err := json.Marshal(order)
		require.NoError(t, err)

		msg := &sqs.SendMessageInput{
			MessageBody: aws.String(string(b)),
			QueueUrl:    aws.String(url),
		}

		trader.EXPECT().Trade(trade.TradeType(tradeType), amount).Return(order, nil)
		sqsClient.EXPECT().SendMessage(msg).Return(nil, nil)

		s, err := service.New(amount, url, sqsClient, trader)
		require.NoError(t, err)

		err = s.Trade(context.Background(), tradeType)
		require.NoError(t, err)
	})
}