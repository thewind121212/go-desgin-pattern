package main

import (
	"log"
	"runtime"
	"time"
)

// T is a type alias to accept any type.
type T = interface{}

// WorkerPool is a contract for Worker Pool implementation
type WorkerPool interface {
	Run()
	AddTask(task func())
	NewWorkerPool(maxWorker int) WorkerPool
}

type workerPool struct {
	maxWorker   int
	queuedTaskC chan func()
}

func (wp *workerPool) Run() {
	for i := 0; i < wp.maxWorker; i++ {
		go func(workerID int) {
			for task := range wp.queuedTaskC {
				task()
			}
		}(i + 1)
	}
}

func (wp *workerPool) AddTask(task func()) {
	wp.queuedTaskC <- task
}

func (wp *workerPool) NewWorkerPool(maxWorker int) WorkerPool {
	return &workerPool{
		maxWorker:   maxWorker,
		queuedTaskC: make(chan func()),
	}
}

func main() {
	// For monitoring purpose.
	waitC := make(chan bool)
	go func() {
		for {
			log.Printf("[main] Total current goroutine: %d", runtime.NumGoroutine())
			time.Sleep(1 * time.Second)
		}
	}()

	// Start Worker Pool.
	totalWorker := 3
	var wp WorkerPool = &workerPool{}
	wp = wp.NewWorkerPool(totalWorker)
	wp.Run()

	type result struct {
		id    int
		value int
	}

	totalTask := 5
	resultC := make(chan result, totalTask)

	for i := 0; i < totalTask; i++ {
		id := i + 1
		wp.AddTask(func() {
			log.Printf("[main] Starting task %d", id)
			time.Sleep(5 * time.Second)
			resultC <- result{id, id * 2}
		})
	}

	for i := 0; i < totalTask; i++ {
		res := <-resultC
		log.Printf("[main] Task %d has been finished with result %d:", res.id, res.value)
	}

	<-waitC
}
