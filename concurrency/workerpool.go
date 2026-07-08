package main

import (
	"fmt"
	"sync"
)

// --- Worker Pool Pattern ---
//
// A fixed pool of N goroutines reads from a shared jobs channel.
// The pool size caps concurrency, preventing resource exhaustion
// (e.g., too many open DB connections or file handles).
//
// Unlike fan-out (which fans a channel into N dedicated goroutines
// and back out again), a worker pool is a long-lived set of
// goroutines that each loop over the same jobs channel until it
// is closed.
//
//  jobs ──► [worker 0]
//        ──► [worker 1]  ──► results
//        ──► [worker 2]

// Job holds the input for a single unit of work.
type Job struct {
	ID    int
	Value int
}

// Result holds the output produced by a worker.
type Result struct {
	JobID  int
	Output int
	Worker int
}

// processJob is the computation each worker performs.
func processJob(job Job, workerID int) Result {
	return Result{
		JobID:  job.ID,
		Output: job.Value * job.Value,
		Worker: workerID,
	}
}

// startWorkerPool creates 'numWorkers' goroutines.
// Each goroutine loops over 'jobs' until the channel is closed,
// then decrements the WaitGroup. When all workers finish,
// the results channel is closed so the consumer can range over it.
func startWorkerPool(numWorkers int, jobs <-chan Job) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for id := 0; id < numWorkers; id++ {
		workerID := id // capture loop variable before goroutine launch
		go func() {
			defer wg.Done()
			// Each worker loops until 'jobs' is closed.
			for job := range jobs {
				results <- processJob(job, workerID)
			}
		}()
	}

	// Close results once every worker has finished.
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// DemoWorkerPool submits 10 jobs to a pool of 3 workers and
// prints each result as it arrives.
// The caller closes 'jobs' to signal that no more work is coming;
// workers drain the channel and exit naturally.
func DemoWorkerPool() {
	fmt.Println("=== Worker Pool ===")

	const numJobs = 10
	const numWorkers = 3

	// Buffer the jobs channel to the total job count so the
	// producer never blocks while filling it.
	jobs := make(chan Job, numJobs)
	for i := 1; i <= numJobs; i++ {
		jobs <- Job{ID: i, Value: i}
	}
	close(jobs) // tell workers: no more jobs after this

	results := startWorkerPool(numWorkers, jobs)

	for r := range results {
		fmt.Printf("job %2d → %3d  (worker %d)\n", r.JobID, r.Output, r.Worker)
	}
}
