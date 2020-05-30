package aws_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cshep4/kripto/services/data-storer/internal/handler/aws"
	"github.com/cshep4/kripto/services/data-storer/internal/mocks/idempotency"
	"github.com/cshep4/kripto/services/data-storer/internal/mocks/service"
	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"github.com/cshep4/kripto/shared/go/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_StoreTrade(t *testing.T) {
	t.Run("returns error if no sqs messages", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx = log.WithServiceName(context.Background(), log.New("debug"), "test")
		)

		err := handler.StoreTrade(ctx, events.SQSEvent{})
		require.Error(t, err)

		assert.Equal(t, "no sqs message passed to function", err.Error())
	})

	t.Run("does not store trade if msg body invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx   = log.WithServiceName(context.Background(), log.New("debug"), "test")
			event = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: "invalid",
				}},
			}
		)

		idempotencer.EXPECT().Check(gomock.Any(), gomock.Any()).Times(0)
		service.EXPECT().StoreTrade(gomock.Any(), gomock.Any()).Times(0)

		err := handler.StoreTrade(ctx, event)
		require.NoError(t, err)
	})

	t.Run("does not store trade if msg body not valid trade", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx   = log.WithServiceName(context.Background(), log.New("debug"), "test")
			event = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: `{
						"id": "tradeId",
						"funds": "invalid"
					}`,
				}},
			}
		)

		idempotencer.EXPECT().Check(gomock.Any(), gomock.Any()).Times(0)
		service.EXPECT().StoreTrade(gomock.Any(), gomock.Any()).Times(0)

		err := handler.StoreTrade(ctx, event)
		require.NoError(t, err)
	})

	t.Run("does not store trade if error checking idempotency", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx   = log.WithServiceName(context.Background(), log.New("debug"), "test")
			event = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: `{
						"id": "tradeId",
						"side": "buy",
						"productId": "productId",
						"funds": "1",
						"fillFees": "2",
						"filledSize": "3",
						"executedValue": "4"
					}`,
				}},
			}
			testErr = errors.New("error")
		)

		idempotencer.EXPECT().Check(ctx, "tradeId").Return(false, testErr)
		service.EXPECT().StoreTrade(gomock.Any(), gomock.Any()).Times(0)

		err := handler.StoreTrade(ctx, event)
		require.NoError(t, err)
	})

	t.Run("does not store trade if idempotency key exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx   = log.WithServiceName(context.Background(), log.New("debug"), "test")
			event = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: `{
						"id": "tradeId",
						"side": "buy",
						"productId": "productId",
						"funds": "1",
						"fillFees": "2",
						"filledSize": "3",
						"executedValue": "4"
					}`,
				}},
			}
		)

		idempotencer.EXPECT().Check(ctx, "tradeId").Return(true, nil)
		service.EXPECT().StoreTrade(gomock.Any(), gomock.Any()).Times(0)

		err := handler.StoreTrade(ctx, event)
		require.NoError(t, err)
	})

	t.Run("returns error if error storing trade", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx   = log.WithServiceName(context.Background(), log.New("debug"), "test")
			event = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: `{
						"id": "tradeId",
						"side": "buy",
						"productId": "productId",
						"funds": "1",
						"fillFees": "2",
						"filledSize": "3",
						"executedValue": "4"
					}`,
				}},
			}
			testErr = errors.New("error")
			trade   = model.Trade{
				Id:         "tradeId",
				TradeType:  model.Buy,
				ProductId:  "productId",
				SpentFunds: float64(1),
				Fees:       float64(2),
				Value: model.Value{
					BTC: float64(3),
					GBP: float64(4),
				},
			}
		)

		idempotencer.EXPECT().Check(ctx, "tradeId").Return(false, nil)
		service.EXPECT().StoreTrade(ctx, trade).Return(testErr)

		err := handler.StoreTrade(ctx, event)
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns nil if trade stored successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx   = log.WithServiceName(context.Background(), log.New("debug"), "test")
			event = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: `{
						"id": "tradeId",
						"side": "buy",
						"productId": "productId",
						"funds": "1",
						"fillFees": "2",
						"filledSize": "3",
						"executedValue": "4"
					}`,
				}},
			}
			trade = model.Trade{
				Id:         "tradeId",
				TradeType:  model.Buy,
				ProductId:  "productId",
				SpentFunds: float64(1),
				Fees:       float64(2),
				Value: model.Value{
					BTC: float64(3),
					GBP: float64(4),
				},
			}
		)

		idempotencer.EXPECT().Check(ctx, "tradeId").Return(false, nil)
		service.EXPECT().StoreTrade(ctx, trade).Return(nil)

		err := handler.StoreTrade(ctx, event)
		require.NoError(t, err)
	})
}

func TestHandler_StoreRate(t *testing.T) {
	t.Run("returns error if no sqs messages", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx = context.Background()
		)

		err := handler.StoreRate(ctx, events.SQSEvent{})
		require.Error(t, err)

		assert.Equal(t, "no sqs message passed to function", err.Error())
	})

	t.Run("does not store trade if msg body invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx   = context.Background()
			event = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: "invalid",
				}},
			}
		)

		idempotencer.EXPECT().Check(gomock.Any(), gomock.Any()).Times(0)
		service.EXPECT().StoreRate(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		err := handler.StoreRate(ctx, event)
		require.NoError(t, err)
	})

	t.Run("does not store trade if error checking idempotency", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx            = log.WithServiceName(context.Background(), log.New("debug"), "test")
			idempotencyKey = "idempotency key"
			now            = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			event          = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: fmt.Sprintf(`{
						"idempotencyKey": "%s",
						"rate": 123.45,
						"dateTime": "%s"
					}`, idempotencyKey, now),
				}},
			}
			testErr = errors.New("error")
		)

		idempotencer.EXPECT().Check(ctx, idempotencyKey).Return(false, testErr)
		service.EXPECT().StoreRate(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		err := handler.StoreRate(ctx, event)
		require.NoError(t, err)
	})

	t.Run("does not store trade if idempotency key exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx            = log.WithServiceName(context.Background(), log.New("debug"), "test")
			idempotencyKey = "idempotency key"
			now            = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			event          = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: fmt.Sprintf(`{
						"idempotencyKey": "%s",
						"rate": 123.45,
						"dateTime": "%s"
					}`, idempotencyKey, now),
				}},
			}
		)

		idempotencer.EXPECT().Check(ctx, idempotencyKey).Return(true, nil)
		service.EXPECT().StoreRate(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		err := handler.StoreRate(ctx, event)
		require.NoError(t, err)
	})

	t.Run("returns error if error storing trade", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx            = log.WithServiceName(context.Background(), log.New("debug"), "test")
			idempotencyKey = "idempotency key"
			now            = time.Now().UTC().Round(time.Second)
			rate           = 123.45
			event          = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: fmt.Sprintf(`{
						"idempotencyKey": "%s",
						"rate": %v,
						"dateTime": "%s"
					}`, idempotencyKey, rate, now.Format("2006-01-02T15:04:05Z")),
				}},
			}
			testErr = errors.New("error")
		)

		idempotencer.EXPECT().Check(ctx, idempotencyKey).Return(false, nil)
		service.EXPECT().StoreRate(ctx, rate, now).Return(testErr)

		err := handler.StoreRate(ctx, event)
		require.Error(t, err)

		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("returns nil if trade stored successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			service      = service_mocks.NewMockServicer(ctrl)
			idempotencer = idempotency_mocks.NewMockIdempotencer(ctrl)
			handler      = aws.Handler{
				Service:      service,
				Idempotencer: idempotencer,
			}
			ctx            = log.WithServiceName(context.Background(), log.New("debug"), "test")
			idempotencyKey = "idempotency key"
			now            = time.Now().UTC().Round(time.Second)
			rate           = 123.45
			event          = events.SQSEvent{
				Records: []events.SQSMessage{{
					Body: fmt.Sprintf(`{
						"idempotencyKey": "%s",
						"rate": %v,
						"dateTime": "%s"
					}`, idempotencyKey, rate, now.Format("2006-01-02T15:04:05Z")),
				}},
			}
		)

		idempotencer.EXPECT().Check(ctx, idempotencyKey).Return(false, nil)
		service.EXPECT().StoreRate(ctx, rate, now).Return(nil)

		err := handler.StoreRate(ctx, event)
		require.NoError(t, err)
	})
}
