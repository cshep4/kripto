package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cshep4/kripto/services/trader/internal/mocks/publish"
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

	t.Run("returns error if topic is empty", func(t *testing.T) {
		s, err := service.New("amount", "", nil, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "topic", ipErr.Parameter)
	})

	t.Run("returns error if publisher is empty", func(t *testing.T) {
		s, err := service.New("amount", "url", nil, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "publisher", ipErr.Parameter)
	})

	t.Run("returns error if trader is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		publisher := publish_mocks.NewMockPublisher(ctrl)

		s, err := service.New("amount", "url", publisher, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "trader", ipErr.Parameter)
	})

	t.Run("returns service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		publisher := publish_mocks.NewMockPublisher(ctrl)
		trader := trader_mocks.NewMockTrader(ctrl)

		s, err := service.New("amount", "url", publisher, trader)
		require.NoError(t, err)

		assert.NotNil(t, s)
	})
}

func TestService_Trade(t *testing.T) {
	t.Run("returns error if error trading", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		publisher := publish_mocks.NewMockPublisher(ctrl)
		trader := trader_mocks.NewMockTrader(ctrl)

		const (
			amount    = "amount"
			url       = "url"
			tradeType = "type"
		)
		testErr := errors.New("error")

		trader.EXPECT().Trade(trade.TradeType(tradeType), amount).Return(nil, testErr)

		s, err := service.New(amount, url, publisher, trader)
		require.NoError(t, err)

		err = s.Trade(context.Background(), tradeType)
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns error if error sending message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		publisher := publish_mocks.NewMockPublisher(ctrl)
		trader := trader_mocks.NewMockTrader(ctrl)

		const (
			amount    = "amount"
			topic     = "topic"
			tradeType = "type"
			id        = "tradeId"
		)
		testErr := errors.New("error")

		order := &trade.TradeResponse{
			Id: id,
		}
		b, err := json.Marshal(order)
		require.NoError(t, err)

		publishInput := &sns.PublishInput{
			Message:  aws.String(string(b)),
			TopicArn: aws.String(topic),
		}
		ctx := context.Background()

		trader.EXPECT().Trade(trade.TradeType(tradeType), amount).Return(order, nil)
		publisher.EXPECT().PublishWithContext(ctx, publishInput).Return(nil, testErr)

		s, err := service.New(amount, topic, publisher, trader)
		require.NoError(t, err)

		err = s.Trade(ctx, tradeType)
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns nil if trade executed successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		publisher := publish_mocks.NewMockPublisher(ctrl)
		trader := trader_mocks.NewMockTrader(ctrl)

		const (
			amount    = "amount"
			topic     = "topic"
			tradeType = "type"
			id        = "tradeId"
		)

		order := &trade.TradeResponse{
			Id: id,
		}
		b, err := json.Marshal(order)
		require.NoError(t, err)

		publishInput := &sns.PublishInput{
			Message:  aws.String(string(b)),
			TopicArn: aws.String(topic),
		}
		ctx := context.Background()

		trader.EXPECT().Trade(trade.TradeType(tradeType), amount).Return(order, nil)
		publisher.EXPECT().PublishWithContext(ctx, publishInput).Return(nil, nil)

		s, err := service.New(amount, topic, publisher, trader)
		require.NoError(t, err)

		err = s.Trade(ctx, tradeType)
		require.NoError(t, err)
	})
}
