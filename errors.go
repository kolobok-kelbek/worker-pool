package main

import "errors"

var ErrNotAllTasksAddToPool = errors.New("not all tasks add to pool")
var ErrNotFoundProcessor = errors.New("not found processor for task")
var ErrChannelClosed = errors.New("channel closed")
