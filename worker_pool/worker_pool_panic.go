package worker_pool

import (
	"fmt"
	"log"
)

func processJobSafely(workerID int, job Job) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf(
				"Worker %d panicked on job %d: %v",
				workerID, job.job_id, r,
			)
		}
	}()
	processJob(workerID, job)
}

func processJob(worker_id int, job Job) {
	if job.job_id == 500 {
		panic("something went wrong")
	}

	fmt.Printf("Worker %v, Processing job id%v\n", worker_id, job.job_id)
}
