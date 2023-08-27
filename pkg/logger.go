package pkg

import "context"

type RequestLog struct {
	RequestID string
	Body      string
}

type ResponseLog struct {
	RequestID      string
	ResponseStatus string
}

type LogsService interface {
	AddRequest(ctx context.Context, requestLog *RequestLog) error
	AddResponse(ctx context.Context, responseLog *ResponseLog) error
	Flush() error
}
