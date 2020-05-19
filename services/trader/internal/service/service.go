package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/cshep4/kripto/services/trader/internal/trader"
)

type (
	Trader interface {
		Trade(tradeType trader.TradeType, amount string) (*trader.TradeResponse, error)
	}

	SQS interface {
		SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	}

	service struct {
		amount    string
		queueUrl  string
		sqsClient SQS
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

func New(amount, queueUrl string, sqsClient SQS, trader Trader) (*service, error) {
	switch {
	case amount == "":
		return nil, InvalidParameterError{Parameter: "amount"}
	case queueUrl == "":
		return nil, InvalidParameterError{Parameter: "queueUrl"}
	case sqsClient == nil:
		return nil, InvalidParameterError{Parameter: "sqsClient"}
	case trader == nil:
		return nil, InvalidParameterError{Parameter: "trader"}
	}

	return &service{
		amount:    amount,
		queueUrl:  queueUrl,
		sqsClient: sqsClient,
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

	_, err = s.sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(b)),
		QueueUrl:    aws.String(s.queueUrl),
	})
	if err != nil {
		return fmt.Errorf("send_message: %w", err)
	}

	return nil
}
