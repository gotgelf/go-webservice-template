package pkg

import "context"

type Log struct {
	RequestID          string
	Body               string
	ResponseStatusCode int
}

type LogsService interface {
	addLog(ctx context.Context, requestLog *Log) error
	Flush() error
}
