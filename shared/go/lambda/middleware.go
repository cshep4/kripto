package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type (
	preExecutorFunc  func(ctx context.Context, payload []byte) (context.Context, []byte, error)
	postExecutorFunc func(ctx context.Context, payload []byte) error

	preExecutor struct {
		before  preExecutorFunc
		handler lambda.Handler
	}

	postExecutor struct {
		after   postExecutorFunc
		handler lambda.Handler
	}
)

func (pe *preExecutor) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	ctx, payload, err := pe.before(ctx, payload)
	if err != nil {
		return nil, err
	}

	return pe.handler.Invoke(ctx, payload)
}

func (pe *postExecutor) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	res, err := pe.handler.Invoke(ctx, payload)
	if err != nil {
		return res, err
	}

	err = pe.after(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}
