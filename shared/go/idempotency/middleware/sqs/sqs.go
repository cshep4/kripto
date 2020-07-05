package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cshep4/kripto/shared/go/idempotency"
	"github.com/cshep4/kripto/shared/go/log"
)

type (
	middleware struct {
		idempotencer idempotency.Idempotencer
	}

	// InvalidParameterError is returned when a required parameter passed to NewMiddleware is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func NewMiddleware(idempotencer idempotency.Idempotencer) (*middleware, error) {
	if idempotencer == nil {
		return nil, InvalidParameterError{Parameter: "idempotencer"}
	}

	return &middleware{
		idempotencer: idempotencer,
	}, nil
}

func (m *middleware) PreExecute(ctx context.Context, payload []byte) (bool, context.Context, []byte, error) {
	var sqsEvent events.SQSEvent
	err := json.Unmarshal(payload, &sqsEvent)
	if err != nil {
		return false, nil, nil, fmt.Errorf("unmarshal: %w", err)
	}

	messages := make([]events.SQSMessage, 0, len(sqsEvent.Records))

	for _, msg := range sqsEvent.Records {
		res, err := m.idempotencer.Check(ctx, msg.MessageId)
		if err != nil {
			log.Error(ctx, "error_checking_idempotency",
				log.SafeParam("id", msg.MessageId),
				log.SafeParam("body", msg.Body),
				log.ErrorParam(err),
			)
			return false, nil, nil, fmt.Errorf("check: %w", err)
		}

		if res.Exists {
			log.Info(ctx, "msg_already_processed",
				log.SafeParam("id", msg.MessageId),
				log.SafeParam("body", msg.Body),
			)
			continue
		}

		err = m.idempotencer.MarkComplete(ctx, msg.MessageId, nil)
		if err != nil {
			log.Error(ctx, "error_marking_complete",
				log.SafeParam("id", msg.MessageId),
				log.SafeParam("body", msg.Body),
				log.ErrorParam(err),
			)
			return false, nil, nil, fmt.Errorf("check: %w", err)
		}

		messages = append(messages, msg)
	}
	sqsEvent.Records = messages

	payload, err = json.Marshal(sqsEvent)
	if err != nil {
		return false, nil, nil, fmt.Errorf("marshal: %w", err)
	}

	return false, ctx, payload, nil
}

func (m *middleware) PostExecute(ctx context.Context, payload, response []byte) error     { return nil }
func (m *middleware) HandleError(ctx context.Context, payload []byte, err error) {}
