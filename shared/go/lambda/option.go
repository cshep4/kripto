package lambda

type option func(*Runner)

func WithPreExecute(e preExecutorFunc) option {
	return func(r *Runner) {
		r.Handler = &preExecutor{
			before:  e,
			handler: r.Handler,
		}
	}
}

func WithPostExecute(e postExecutorFunc) option {
	return func(r *Runner) {
		r.Handler = &postExecutor{
			after:   e,
			handler: r.Handler,
		}
	}
}
