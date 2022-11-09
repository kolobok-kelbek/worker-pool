package main

type Result struct {
	Task  Task
	Data  interface{}
	Error error
}

func NewResult(task Task, data interface{}, err error) Result {
	return Result{
		Task:  task,
		Data:  data,
		Error: err,
	}
}

func NewSuccessResult(task Task, data interface{}) Result {
	return NewResult(task, data, nil)
}

func NewFailResult(task Task, err error) Result {
	return NewResult(task, nil, err)
}
