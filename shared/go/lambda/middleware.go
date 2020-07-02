package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type (
	preExecutorFunc  func(ctx context.Context, payload []byte) (bool, context.Context, []byte, error)
	postExecutorFunc func(ctx context.Context, payload []byte) error
	errorHandlerFunc func(ctx context.Context, err error)

	preExecutor struct {
		runner  *runner
		before  preExecutorFunc
		handler lambda.Handler
	}

	postExecutor struct {
		runner  *runner
		after   postExecutorFunc
		handler lambda.Handler
	}

	errorHandler struct {
		errorHandler errorHandlerFunc
		handler      lambda.Handler
	}
)

func (pe *preExecutor) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	done, ctx, payload, err := pe.before(ctx, payload)
	if err != nil {
		return nil, err
	}
	if done {
		pe.runner.terminated = true
		return payload, nil
	}

	return pe.handler.Invoke(ctx, payload)
}

func (pe *postExecutor) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	res, err := pe.handler.Invoke(ctx, payload)
	if err != nil {
		return nil, err
	}
	if pe.runner.terminated {
		return res, nil
	}

	err = pe.after(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (eh *errorHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	res, err := eh.handler.Invoke(ctx, payload)
	if err == nil {
		return res, nil
	}

	eh.errorHandler(ctx, err)

	return nil, err
}
