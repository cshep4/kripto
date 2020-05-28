package model_test

import (
	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTradeRequest_ToTrade(t *testing.T) {
	t.Run("return error if id is empty", func(t *testing.T) {
		trade, err := (&model.TradeRequest{}).ToTrade()
		require.Error(t, err)

		assert.Empty(t, trade)

		ipErr, ok := err.(model.InvalidPropertyError)
		assert.True(t, ok)
		assert.Equal(t, "id", ipErr.Parameter)
		assert.Equal(t, "value is empty", ipErr.Err)
	})

	t.Run("return error if side is empty", func(t *testing.T) {
		trade, err := (&model.TradeRequest{
			Id: "id",
		}).ToTrade()
		require.Error(t, err)

		assert.Empty(t, trade)

		ipErr, ok := err.(model.InvalidPropertyError)
		assert.True(t, ok)
		assert.Equal(t, "side", ipErr.Parameter)
		assert.Equal(t, "value is empty", ipErr.Err)
	})

	t.Run("return error if productId is empty", func(t *testing.T) {
		trade, err := (&model.TradeRequest{
			Id:   "id",
			Side: "side",
		}).ToTrade()
		require.Error(t, err)

		assert.Empty(t, trade)

		ipErr, ok := err.(model.InvalidPropertyError)
		assert.True(t, ok)
		assert.Equal(t, "productId", ipErr.Parameter)
		assert.Equal(t, "value is empty", ipErr.Err)
	})

	t.Run("return error if funds is invalid", func(t *testing.T) {
		trade, err := (&model.TradeRequest{
			Id:        "id",
			Side:      "side",
			ProductId: "productId",
			Funds:     "invalid",
		}).ToTrade()
		require.Error(t, err)

		assert.Empty(t, trade)

		ipErr, ok := err.(model.InvalidPropertyError)
		assert.True(t, ok)
		assert.Equal(t, "funds", ipErr.Parameter)
		assert.Equal(t, "strconv.ParseFloat: parsing \"invalid\": invalid syntax", ipErr.Err)
	})

	t.Run("return error if fillFees is invalid", func(t *testing.T) {
		trade, err := (&model.TradeRequest{
			Id:        "id",
			Side:      "side",
			ProductId: "productId",
			Funds:     "1",
			FillFees:  "invalid",
		}).ToTrade()
		require.Error(t, err)

		assert.Empty(t, trade)

		ipErr, ok := err.(model.InvalidPropertyError)
		assert.True(t, ok)
		assert.Equal(t, "fillFees", ipErr.Parameter)
		assert.Equal(t, "strconv.ParseFloat: parsing \"invalid\": invalid syntax", ipErr.Err)
	})

	t.Run("return error if filledSize is invalid", func(t *testing.T) {
		trade, err := (&model.TradeRequest{
			Id:         "id",
			Side:       "side",
			ProductId:  "productId",
			Funds:      "1",
			FillFees:   "2",
			FilledSize: "invalid",
		}).ToTrade()
		require.Error(t, err)

		assert.Empty(t, trade)

		ipErr, ok := err.(model.InvalidPropertyError)
		assert.True(t, ok)
		assert.Equal(t, "filledSize", ipErr.Parameter)
		assert.Equal(t, "strconv.ParseFloat: parsing \"invalid\": invalid syntax", ipErr.Err)
	})

	t.Run("return error if executedValue is invalid", func(t *testing.T) {
		trade, err := (&model.TradeRequest{
			Id:            "id",
			Side:          "side",
			ProductId:     "productId",
			Funds:         "1",
			FillFees:      "2",
			FilledSize:    "3",
			ExecutedValue: "invalid",
		}).ToTrade()
		require.Error(t, err)

		assert.Empty(t, trade)

		ipErr, ok := err.(model.InvalidPropertyError)
		assert.True(t, ok)
		assert.Equal(t, "executedValue", ipErr.Parameter)
		assert.Equal(t, "strconv.ParseFloat: parsing \"invalid\": invalid syntax", ipErr.Err)
	})

	t.Run("returns trade", func(t *testing.T) {
		now := time.Now()
		req := model.TradeRequest{
			Id:            "id",
			Side:          "buy",
			ProductId:     "productId",
			Settled:       true,
			CreatedAt:     model.Time{Time: now},
			Funds:         "1",
			FillFees:      "2",
			FilledSize:    "3",
			ExecutedValue: "4",
		}

		trade, err := req.ToTrade()
		require.NoError(t, err)

		assert.Equal(t, req.Id, trade.Id)
		assert.Equal(t, req.Side, trade.TradeType)
		assert.Equal(t, req.ProductId, trade.ProductId)
		assert.Equal(t, req.Settled, trade.Settled)
		assert.Equal(t, now, trade.CreatedAt.Time)
		assert.Equal(t, float64(1), trade.SpentFunds)
		assert.Equal(t, float64(2), trade.Fees)
		assert.Equal(t, float64(3), trade.Value.BTC)
		assert.Equal(t, float64(4), trade.Value.GBP)
	})
}
