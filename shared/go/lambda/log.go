package lambda

import (
	"context"
	"github.com/cshep4/kripto/shared/go/log"
)

func LogMiddleware(level, service, function string) preExecutorFunc {
	return func(ctx context.Context, payload []byte) (context.Context, []byte, error) {
		ctx = log.WithFunctionName(ctx,
			log.New(level),
			service,
			function,
		)

		return ctx, payload, nil
	}
}
