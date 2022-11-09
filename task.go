package main

import "context"

type Task struct {
	ID          string
	ProcessName string
	Data        interface{}
	Context     context.Context
}
