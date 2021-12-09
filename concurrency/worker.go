package concurrency

import (
	"sync"
	"sync/atomic"
)

type WorkerPool struct {
	queue     chan job
	wg        sync.WaitGroup
	processed int64
	limit     int64
	running   int32
}

type job struct {
	task func() bool
}

func NewWorkerPool(workerCount int, entriesPerWorker int, limit int) *WorkerPool {
	wp := &WorkerPool{
		queue: make(chan job),
		limit: int64(limit),
	}

	wp.AddWorkers(workerCount, entriesPerWorker)

	return wp
}

func (wp *WorkerPool) AddWorkers(workerCount int, entriesPerWorker int) {
	wp.wg.Add(workerCount)
	defer wp.adjustCount(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer func() {
				wp.adjustCount(-1)
				wp.wg.Done()
			}()
			processed := 0
			for j := range wp.queue {
				if wp.limitReached() || (entriesPerWorker > 0 && processed == entriesPerWorker) {
					continue
				}
				result := j.task()
				if result {
					atomic.AddInt64(&wp.processed, 1)
					processed++
				}
			}
		}()
	}
}

func (wp *WorkerPool) limitReached() bool {
	if wp.limit == 0 {
		return false
	}
	return atomic.LoadInt64(&wp.processed) >= wp.limit
}

func (wp *WorkerPool) adjustCount(delta int) {
	atomic.AddInt32(&wp.running, int32(delta))
}

func (wp *WorkerPool) isRunning() bool {
	return atomic.LoadInt32(&wp.running) > 0
}

func (wp *WorkerPool) Push(task func() bool) {
	if !wp.limitReached() && wp.isRunning() {
		select {
		case wp.queue <- job{task: task}:
		}
	}
}

func (wp *WorkerPool) Close() {
	close(wp.queue)
	wp.wg.Wait()
}
