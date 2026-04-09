package mbjobs

import (
	"sync"
)

// Job is a function to be executed in the background.
type Job func()

var (
	jobChan chan Job
	once    sync.Once
)

// EnsureWorkerPool initializes the background worker pool if it hasn't been started.
func EnsureWorkerPool() {
	once.Do(func() {
		// Default to 4 workers for background I/O and processing.
		jobChan = make(chan Job, 1024)
		for i := 0; i < 4; i++ {
			go func() {
				for job := range jobChan {
					job()
				}
			}()
		}
	})
}

// EnqueueJob schedules a task to run on the background worker pool.
func EnqueueJob(j Job) {
	EnsureWorkerPool()
	jobChan <- j
}
