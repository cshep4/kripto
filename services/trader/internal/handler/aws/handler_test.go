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
