package service_test

import (
	"context"
	"errors"
	"github.com/cshep4/kripto/services/data-storer/internal/mocks/rate"
	"github.com/cshep4/kripto/services/data-storer/internal/mocks/trade"
	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"github.com/cshep4/kripto/services/data-storer/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("returns error if rateStore is empty", func(t *testing.T) {
		s, err := service.New(nil, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "rateStore", ipErr.Parameter)
	})

	t.Run("returns error if tradeStore is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rateStore := rate_mocks.NewMockRateStore(ctrl)

		s, err := service.New(rateStore, nil)
		require.Error(t, err)

		assert.Nil(t, s)

		ipErr, ok := err.(service.InvalidParameterError)
		assert.True(t, ok)
		assert.Equal(t, "tradeStore", ipErr.Parameter)
	})

	t.Run("returns service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rateStore := rate_mocks.NewMockRateStore(ctrl)
		tradeStore := trade_mocks.NewMockTradeStore(ctrl)

		s, err := service.New(rateStore, tradeStore)
		require.NoError(t, err)

		assert.NotNil(t, s)
	})
}

func TestService_StoreRate(t *testing.T) {
	t.Run("returns error if error storing rate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rateStore := rate_mocks.NewMockRateStore(ctrl)
		tradeStore := trade_mocks.NewMockTradeStore(ctrl)

		var (
			ctx     = context.Background()
			now     = time.Now()
			testErr = errors.New("error")
		)
		const rate = float64(12.34)

		s, err := service.New(rateStore, tradeStore)
		require.NoError(t, err)

		rateStore.EXPECT().Store(ctx, rate, now).Return(testErr)

		err = s.StoreRate(ctx, rate, now)
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns nil if rate stored successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rateStore := rate_mocks.NewMockRateStore(ctrl)
		tradeStore := trade_mocks.NewMockTradeStore(ctrl)

		var (
			ctx = context.Background()
			now = time.Now()
		)
		const rate = float64(12.34)

		s, err := service.New(rateStore, tradeStore)
		require.NoError(t, err)

		rateStore.EXPECT().Store(ctx, rate, now).Return(nil)

		err = s.StoreRate(ctx, rate, now)
		require.NoError(t, err)
	})
}

func TestService_StoreTrade(t *testing.T) {
	t.Run("returns error if error storing trade", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			rateStore  = rate_mocks.NewMockRateStore(ctrl)
			tradeStore = trade_mocks.NewMockTradeStore(ctrl)

			ctx     = context.Background()
			testErr = errors.New("error")
			trade   = model.Trade{
				Id: "id",
			}
		)

		s, err := service.New(rateStore, tradeStore)
		require.NoError(t, err)

		tradeStore.EXPECT().Store(ctx, trade).Return(testErr)

		err = s.StoreTrade(ctx, trade)
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns nil if trade stored successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			rateStore  = rate_mocks.NewMockRateStore(ctrl)
			tradeStore = trade_mocks.NewMockTradeStore(ctrl)

			ctx   = context.Background()
			trade = model.Trade{
				Id: "id",
			}
		)

		s, err := service.New(rateStore, tradeStore)
		require.NoError(t, err)

		tradeStore.EXPECT().Store(ctx, trade).Return(nil)

		err = s.StoreTrade(ctx, trade)
		require.NoError(t, err)
	})
}
