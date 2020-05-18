package trader_test

import (
	"errors"
	"testing"

	"github.com/cshep4/kripto/services/trader/internal/mocks/coinbase"
	"github.com/cshep4/kripto/services/trader/internal/service"
	"github.com/cshep4/kripto/services/trader/internal/trader"
	"github.com/golang/mock/gomock"
	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("returns error if coinbase is empty", func(t *testing.T) {
		trader, err := trader.New(nil)
		require.Error(t, err)

		assert.Nil(t, trader)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "coinbase", ipErr.Parameter)
	})

	t.Run("returns service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		coinbase := coinbase_mocks.NewMockCoinbase(ctrl)

		trader, err := trader.New(coinbase)
		require.NoError(t, err)

		assert.NotNil(t, trader)
	})
}

func TestTrader_Trade(t *testing.T) {
	t.Run("returns error if error creating order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const (
			amount                     = "amount"
			tradeType trader.TradeType = "tradeType"
			productId                  = "BTC-GBP"
			orderType                  = "market"
		)

		coinbase := coinbase_mocks.NewMockCoinbase(ctrl)

		trader, err := trader.New(coinbase)
		require.NoError(t, err)

		order := &coinbasepro.Order{
			Funds:     amount,
			Side:      string(tradeType),
			ProductID: productId,
			Type:      orderType,
		}
		testErr := errors.New("error")

		coinbase.EXPECT().CreateOrder(order).Return(coinbasepro.Order{}, testErr)

		res, err := trader.Trade(tradeType, amount)
		require.Error(t, err)

		assert.Empty(t, res)
	})

	t.Run("returns error if error getting order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const (
			amount                     = "amount"
			tradeType trader.TradeType = "tradeType"
			productId                  = "BTC-GBP"
			orderType                  = "market"
			orderId                    = "id"
		)

		coinbase := coinbase_mocks.NewMockCoinbase(ctrl)

		trader, err := trader.New(coinbase)
		require.NoError(t, err)

		order := &coinbasepro.Order{
			Funds:     amount,
			Side:      string(tradeType),
			ProductID: productId,
			Type:      orderType,
		}
		orderRes := coinbasepro.Order{
			ID: orderId,
		}
		testErr := errors.New("error")

		coinbase.EXPECT().CreateOrder(order).Return(orderRes, nil)
		coinbase.EXPECT().GetOrder(orderId).Return(coinbasepro.Order{}, testErr)

		res, err := trader.Trade(tradeType, amount)
		require.Error(t, err)

		assert.Empty(t, res)
	})

	t.Run("returns trade response if executed successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const (
			amount                     = "amount"
			tradeType trader.TradeType = "tradeType"
			productId                  = "BTC-GBP"
			orderType                  = "market"
			orderId                    = "id"
		)

		coinbase := coinbase_mocks.NewMockCoinbase(ctrl)

		trader, err := trader.New(coinbase)
		require.NoError(t, err)

		order := &coinbasepro.Order{
			Funds:     amount,
			Side:      string(tradeType),
			ProductID: productId,
			Type:      orderType,
		}
		orderRes := coinbasepro.Order{
			ID: orderId,
		}

		coinbase.EXPECT().CreateOrder(order).Return(orderRes, nil)
		coinbase.EXPECT().GetOrder(orderId).Return(orderRes, nil)

		res, err := trader.Trade(tradeType, amount)
		require.NoError(t, err)

		assert.Equal(t, orderRes.ID, res.ID)
	})
}
