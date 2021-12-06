package concurrency

import (
	"sync"
)

type WorkerPool struct {
	queue chan job
	wg    sync.WaitGroup
}

type job struct {
	task func() bool
}

func NewWorkerPool(workerCount int, entriesPerWorker int) *WorkerPool {
	wp := &WorkerPool{
		queue: make(chan job),
	}

	wp.AddWorkers(workerCount, entriesPerWorker)

	return wp
}

func (wp *WorkerPool) AddWorkers(workerCount int, entriesPerWorker int) {
	wp.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wp.wg.Done()
			processed := 0
			for j := range wp.queue {
				result := j.task()
				if result {
					processed++
				}

				if entriesPerWorker > 0 && processed == entriesPerWorker {
					break
				}
			}
		}()
	}
}

func (wp *WorkerPool) Push(task func() bool) {
	select {
	case wp.queue <- job{task: task}:
	}
}

func (wp *WorkerPool) Close() {
	close(wp.queue)
	wp.wg.Wait()
}
