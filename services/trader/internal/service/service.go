package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cshep4/kripto/services/trader/internal/trader"
)

type (
	Trader interface {
		Trade(tradeType trader.TradeType, amount string) (*trader.TradeResponse, error)
	}
	Publisher interface {
		PublishWithContext(ctx context.Context, input *sns.PublishInput, opts ...request.Option) (*sns.PublishOutput, error)
	}

	service struct {
		amount    string
		topic     string
		publisher Publisher
		trader    Trader
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func New(amount, topic string, publisher Publisher, trader Trader) (*service, error) {
	switch {
	case amount == "":
		return nil, InvalidParameterError{Parameter: "amount"}
	case topic == "":
		return nil, InvalidParameterError{Parameter: "topic"}
	case publisher == nil:
		return nil, InvalidParameterError{Parameter: "publisher"}
	case trader == nil:
		return nil, InvalidParameterError{Parameter: "trader"}
	}

	return &service{
		amount:    amount,
		topic:     topic,
		publisher: publisher,
		trader:    trader,
	}, nil
}

func (s *service) Trade(ctx context.Context, tradeType string) error {
	order, err := s.trader.Trade(trader.TradeType(tradeType), s.amount)
	if err != nil {
		return fmt.Errorf("trade: %w", err)
	}

	b, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("json_marshal: %w", err)
	}

	_, err = s.publisher.PublishWithContext(ctx, &sns.PublishInput{
		Message:  aws.String(string(b)),
		TopicArn: aws.String(s.topic),
	})
	if err != nil {
		return fmt.Errorf("send_message: %w", err)
	}

	return nil
}
