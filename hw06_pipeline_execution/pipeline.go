package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		if done != nil {
			out = func(in In, done In) Out {
				out := make(Bi)
				go func() {
					defer close(out)
					for {
						select {
						case v, ok := <-in:
							if !ok {
								return
							}
							out <- v
						case <-done:
							return
						}
					}
				}()
				return out
			}(out, done)
		}
		out = stage(out)
	}
	return out
}
