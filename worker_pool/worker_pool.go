package worker_pool

import (
	"sync"
)

// Implement a worker pool with 5 workers
// Problem: Given 1,000 jobs, process them using five workers.

type Job struct {
	job_id int
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		processJobSafely(id, job)

	}
}

func main_worker_pool() {
	var wg sync.WaitGroup

	jobs := make(chan Job)

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	for i := 1; i <= 1000; i++ {
		jobs <- Job{job_id: i}
	}

	close(jobs)
	wg.Wait()

}
