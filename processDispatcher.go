package main

import "sync"

type ProcessDispatcherImpl struct {
	processors sync.Map
}

func NewProcessDispatcher() *ProcessDispatcherImpl {
	return &ProcessDispatcherImpl{}
}

func (dispatcher *ProcessDispatcherImpl) RegisterProcessor(processor Processor) {
	dispatcher.processors.Store(processor.GetProcessName(), processor)
}

func (dispatcher *ProcessDispatcherImpl) GetProcessor(processName string) (processor Processor, has bool) {
	var value interface{}
	value, has = dispatcher.processors.Load(processName)

	if !has {
		return nil, false
	}

	processor, has = value.(Processor)
	if !has {
		return nil, false
	}

	return processor, true
}
