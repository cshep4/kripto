package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	uuid "github.com/kevinburke/go.uuid"
)

const requestId value = "requestId"

type (
	value string

	initialiser struct {
		handler lambda.Handler
		runner  *runner
	}
)

func (i *initialiser) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	res, err := i.handler.Invoke(context.WithValue(ctx, requestId, uuid.NewV4().String()), payload)

	go func() {
		i.runner.lock.Lock()
		defer i.runner.lock.Unlock()
		delete(i.runner.terminatedRequests, string(requestId))
	}()

	return res, err
}
