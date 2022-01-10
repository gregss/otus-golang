package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	wg.Add(n)
	cht := make(chan Task)
	var curm int32

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				t, ok := <-cht
				if !ok {
					break
				}
				if atomic.LoadInt32(&curm) < int32(m) || m < 1 {
					err := t()
					if err != nil {
						atomic.AddInt32(&curm, 1)
					}
				}
			}
		}()
	}

	for _, task := range tasks {
		cht <- task
	}
	close(cht)
	wg.Wait()

	if atomic.LoadInt32(&curm) >= int32(m) && m > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
