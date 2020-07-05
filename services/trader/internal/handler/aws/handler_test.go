package aws_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cshep4/kripto/services/trader/internal/handler/aws"
	"github.com/cshep4/kripto/services/trader/internal/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Trade(t *testing.T) {
	t.Run("returns error if tradeType is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx     = context.Background()
			handler = aws.Handler{}
		)

		err := handler.Trade(ctx, aws.TradeRequest{})
		require.Error(t, err)

		brErr, ok := err.(aws.BadRequestError)
		assert.True(t, ok)
		assert.Equal(t, "tradeType", brErr.Parameter)
		assert.Equal(t, "empty", brErr.Err)
	})

	t.Run("returns error if tradeType is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx     = context.Background()
			handler = aws.Handler{}
		)

		err := handler.Trade(ctx, aws.TradeRequest{TradeType: "invalid"})
		require.Error(t, err)

		brErr, ok := err.(aws.BadRequestError)
		assert.True(t, ok)
		assert.Equal(t, "tradeType", brErr.Parameter)
		assert.Equal(t, "invalid value - should be either buy/sell", brErr.Err)
	})

	t.Run("returns error if amount is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx     = context.Background()
			handler = aws.Handler{}
		)

		err := handler.Trade(ctx, aws.TradeRequest{TradeType: "buy"})
		require.Error(t, err)

		brErr, ok := err.(aws.BadRequestError)
		assert.True(t, ok)
		assert.Equal(t, "amount", brErr.Parameter)
		assert.Equal(t, "empty", brErr.Err)
	})

	t.Run("returns error if error making trade", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx     = context.Background()
			testErr = errors.New("error")
		)
		const (
			tradeType = "buy"
			amount    = "amount"
		)

		service := service_mocks.NewMockServicer(ctrl)
		handler := aws.Handler{Service: service}

		service.EXPECT().Trade(ctx, tradeType, amount).Return(testErr)

		err := handler.Trade(ctx, aws.TradeRequest{
			TradeType: tradeType,
			Amount:    amount,
		})
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns nil if successfully traded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		const (
			tradeType = "buy"
			amount    = "amount"
		)

		service := service_mocks.NewMockServicer(ctrl)
		handler := aws.Handler{Service: service}

		service.EXPECT().Trade(ctx, tradeType, amount).Return(nil)

		err := handler.Trade(ctx, aws.TradeRequest{
			TradeType: tradeType,
			Amount:    amount,
		})
		require.NoError(t, err)
	})
}
