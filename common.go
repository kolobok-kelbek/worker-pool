package main

import (
	"context"
)

const (
	ProcessGatewayMetric                     = "process_gateway"
	ProcessDataProcessingMetric              = "process_data_processing"
	ProcessBulkInserterMetric                = "process_bulk_inserter"
	VacanciesBulkToElasticSearchErrorsMetric = "vacancies_bulk_to_elastic_search_errors"
)

type ResultGenerator func() (Result, bool)

type Processor interface {
	Process(id int, task Task) (result interface{}, err error)
	GetProcessName() string
}

type ProcessDispatcher interface {
	RegisterProcessor(processor Processor)
	GetProcessor(processName string) (processor Processor, has bool)
}

type Worker interface {
	Work(int, Task, chan<- Result)
}

type Pool interface {
	Run(ctx context.Context) (ResultGenerator, context.CancelFunc)
	Add(tasks ...Task) error
}
