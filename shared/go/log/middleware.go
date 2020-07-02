package log

import "context"

func Middleware(level, service, function string) func(ctx context.Context, payload []byte) (bool, context.Context, []byte, error) {
	return func(ctx context.Context, payload []byte) (bool, context.Context, []byte, error) {
		return false,
			WithFunctionName(ctx, New(level), service, function),
			payload,
			nil
	}
}
