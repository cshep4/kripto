package idempotency

import "context"

type Middleware interface {
	PreExecute(ctx context.Context, payload []byte) (bool, context.Context, []byte, error)
	PostExecute(ctx context.Context, payload []byte) error
	HandleError(ctx context.Context, err error)
}
