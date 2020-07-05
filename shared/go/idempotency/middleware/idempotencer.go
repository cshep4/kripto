package idempotency

import "context"

type Middleware interface {
	PreExecute(ctx context.Context, payload []byte) (bool, context.Context, []byte, error)
	PostExecute(ctx context.Context, payload, response []byte) error
	HandleError(ctx context.Context, payload []byte, err error)
}
