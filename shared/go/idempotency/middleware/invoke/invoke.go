package invoke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
	var req map[string]interface{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return true, ctx, nil, fmt.Errorf("unmarshal: %w", err)
	}

	idempotencyKey, ok := req["idempotencyKey"].(string)
	if !ok {
		return true, ctx, nil, errors.New("invalid idempotency key")
	}

	res, err := m.idempotencer.Check(ctx, idempotencyKey)
	if err != nil {
		log.Error(ctx, "error_checking_idempotency",
			log.SafeParam("idempotencyKey", idempotencyKey),
			log.SafeParam("request", string(payload)),
			log.ErrorParam(err),
		)
		return true, ctx, nil, fmt.Errorf("check: %w", err)
	}

	if !res.Exists {
		return false, ctx, payload, nil
	}

	log.Info(ctx, "request_already_processed",
		log.SafeParam("idempotencyKey", idempotencyKey),
		log.SafeParam("request", string(payload)),
	)

	if res.Err != nil {
		return true, ctx, nil, res.Err
	}

	return true, ctx, res.Response, nil
}

func (m *middleware) PostExecute(ctx context.Context, payload, response []byte) error {
	var req map[string]interface{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		log.Error(ctx, "error_unmarshalling_payload",
			log.SafeParam("payload", string(payload)),
			log.ErrorParam(err),
		)
		return nil
	}

	idempotencyKey, ok := req["idempotencyKey"].(string)
	if !ok {
		log.Error(ctx, "invalid_idempotency_key",
			log.SafeParam("payload", string(payload)),
			log.ErrorParam(err),
		)
		return nil
	}

	err = m.idempotencer.MarkComplete(ctx, idempotencyKey, response)
	if err != nil {
		log.Error(ctx, "error_marking_complete",
			log.SafeParam("idempotencyKey", idempotencyKey),
			log.ErrorParam(err),
		)
		return nil
	}

	return nil
}

func (m *middleware) HandleError(ctx context.Context, payload []byte, err error) {
	var req map[string]interface{}
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Error(ctx, "error_unmarshalling_payload",
			log.SafeParam("payload", string(payload)),
			log.ErrorParam(err),
		)
		return
	}

	idempotencyKey, ok := req["idempotencyKey"].(string)
	if !ok {
		log.Error(ctx, "invalid_idempotency_key",
			log.SafeParam("payload", string(payload)),
			log.ErrorParam(err),
		)
		return
	}

	idempotencyErr := m.idempotencer.MarkError(ctx, idempotencyKey, err)
	if idempotencyErr != nil {
		log.Error(ctx, "error_marking_error",
			log.SafeParam("idempotencyKey", idempotencyKey),
			log.ErrorParam(idempotencyErr),
			log.ErrorParam(err),
		)
	}
}
