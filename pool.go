package main

import (
	"context"
	"sync"

	"go.uber.org/atomic"
)

type PoolImpl struct {
	config          Config
	taskCollector   chan Task
	resultCollector chan Result
	isRun           *atomic.Bool

	logger Logger
	worker Worker
}

func NewPool(config Config, worker Worker, logger Logger) *PoolImpl {
	pool := &PoolImpl{
		config:          config,
		taskCollector:   make(chan Task, config.CollectorSize),
		resultCollector: make(chan Result, config.CollectorSize),
		isRun:           atomic.NewBool(false),
		worker:          worker,
		logger:          logger,
	}

	return pool
}

func (pool *PoolImpl) Run(ctx context.Context) (ResultGenerator, context.CancelFunc) {
	pool.logger.Info("worker pool started run")

	workContext, cancel := context.WithCancel(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(pool.config.Concurrency)

	for i := 1; i <= pool.config.Concurrency; i++ {
		go pool.work(i, ctx, wg)
	}

	// Горутина ответственная за прекращение работы worker-ов при получении сигнала из context
	go func(pool *PoolImpl) {
		<-workContext.Done()

		pool.isRun.Store(false)

		wg.Wait()

		close(pool.taskCollector)
		close(pool.resultCollector)

		pool.logger.Info("worker pool stopped")
	}(pool)

	pool.isRun.Store(true)

	resultGenerator := func() (Result, bool) {
		if !pool.isRun.Load() {
			return Result{}, false
		}

		select {
		case <-workContext.Done():
			return Result{}, false
		case result := <-pool.resultCollector:
			return result, true
		}
	}

	return resultGenerator, cancel
}

func (pool *PoolImpl) Add(tasks ...Task) error {
	if !pool.isRun.Load() {
		return ErrChannelClosed
	}

	for _, task := range tasks {
		if !pool.isRun.Load() {
			return ErrNotAllTasksAddToPool
		}

		pool.taskCollector <- task
	}

	return nil
}

func (pool *PoolImpl) work(id int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer pool.logger.Info("worker %d stopped", id)

	pool.logger.Info("worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-pool.taskCollector:
			pool.worker.Work(id, task, pool.resultCollector)
		}
	}
}
