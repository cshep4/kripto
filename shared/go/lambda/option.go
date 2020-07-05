package lambda

const (
	preExecute optionType = iota + 1
	postExecute
	errorHandle
)

type (
	optionType uint

	option struct {
		optionType
		apply func(*runner)
	}
)

func WithPreExecute(e preExecutorFunc) option {
	return option{
		optionType: preExecute,
		apply: func(r *runner) {
			r.Handler = &preExecutor{
				runner:  r,
				before:  e,
				handler: r.Handler,
			}
		},
	}
}

func WithPostExecute(e postExecutorFunc) option {
	return option{
		optionType: postExecute,
		apply: func(r *runner) {
			r.Handler = &postExecutor{
				runner:  r,
				after:   e,
				handler: r.Handler,
			}
		},
	}
}

func WithErrorHandler(e errorHandlerFunc) option {
	return option{
		optionType: errorHandle,
		apply: func(r *runner) {
			r.Handler = &errorHandler{
				runner:       r,
				errorHandler: e,
				handler:      r.Handler,
			}
		},
	}
}
