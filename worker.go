package main

import (
	"errors"
	"fmt"
)

type WorkerImpl struct {
	processDispatcher ProcessDispatcher
	logger            Logger
}

func NewWorkerImpl(processDispatcher ProcessDispatcher, logger Logger) *WorkerImpl {
	return &WorkerImpl{
		processDispatcher: processDispatcher,
		logger:            logger,
	}
}

func (worker *WorkerImpl) Work(id int, task Task, resultCollector chan<- Result) {
	defer func() {
		if recoveryMessage := recover(); recoveryMessage != nil {
			err := worker.recoverErrorHandle(recoveryMessage)

			errorMessage := "worker %d with task %s in %s processor thrown panic: %s"
			worker.logger.Error(errorMessage, id, task.ID, task.ProcessName, err.Error())

			resultCollector <- NewFailResult(task, err)
		}
	}()

	processor, has := worker.processDispatcher.GetProcessor(task.ProcessName)
	if !has {
		worker.logger.Error("worker %d not found processor with name %s for task %s", id, task.ProcessName, task.ID)
		resultCollector <- NewFailResult(task, ErrNotFoundProcessor)

		return
	}

	worker.logger.Debug("worker %d started task %s in %s processor", id, task.ID, task.ProcessName)
	result, err := processor.Process(id, task)

	if err != nil {
		resultCollector <- NewFailResult(task, err)

		worker.logger.Error("worker %d failed task %s in %s processor: %s", id, task.ID, task.ProcessName, err)

		return
	}

	resultCollector <- NewSuccessResult(task, result)

	worker.logger.Debug("worker %d completed task %s in %s processor", id, task.ID, task.ProcessName)
}

func (worker *WorkerImpl) recoverErrorHandle(recoveryMessage any) error {
	switch v := recoveryMessage.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		worker.logger.Warn("not supported type %T!\n", v)

		return errors.New(fmt.Sprintf("%v", v))
	}
}
