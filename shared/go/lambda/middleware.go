package lambda

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

type (
	preExecutorFunc  func(ctx context.Context, payload []byte) (bool, context.Context, []byte, error)
	postExecutorFunc func(ctx context.Context, payload, response []byte) error
	errorHandlerFunc func(ctx context.Context, payload []byte, err error)

	preExecutor struct {
		handler lambda.Handler
		runner  *runner
		before  preExecutorFunc
	}

	postExecutor struct {
		handler lambda.Handler
		runner  *runner
		after   postExecutorFunc
	}

	errorHandler struct {
		handler      lambda.Handler
		runner       *runner
		errorHandler errorHandlerFunc
	}
)

func (pe *preExecutor) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	done, ctx, payload, err := pe.before(ctx, payload)
	if done {
		pe.runner.finish(ctx)
	}
	if err != nil {
		return nil, err
	}
	if done {
		return payload, nil
	}

	return pe.handler.Invoke(ctx, payload)
}

func (pe *postExecutor) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	res, err := pe.handler.Invoke(ctx, payload)
	if err != nil {
		return nil, err
	}
	if pe.runner.done(ctx) {
		return res, err
	}

	err = pe.after(ctx, payload, res)
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
	if eh.runner.done(ctx) {
		return res, err
	}

	eh.errorHandler(ctx, payload, err)

	return nil, err
}
