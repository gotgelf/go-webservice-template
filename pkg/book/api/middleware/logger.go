package middleware

import (
	"docker/pkg"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sync"
)

type LogsService struct {
	requestsLogsService *requestsLogsService
}

type requestsLogsService struct {
	channel   chan *pkg.RequestLog
	shutdown  chan bool
	waitGroup sync.WaitGroup
}

func NewLogsService() *LogsService {
	service := &LogsService{
		requestsLogsService: &requestsLogsService{
			channel:  make(chan *pkg.RequestLog, 500),
			shutdown: make(chan bool),
		},
	}
	service.requestsLogsService.waitGroup.Add(1)
	go service.listen()

	return service
}

func (ls *LogsService) Handle() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		body := string(c.Request().Body())
		requestID := c.Locals("requestid").(string)
		//response := c.Response().StatusCode()
		_ = ls.addRequestLog(&pkg.RequestLog{
			RequestID: requestID,
			Body:      body,
		})

		return c.Next()
	}
}

func (ls *LogsService) listen() {
	defer ls.requestsLogsService.waitGroup.Done()
	var logs []*pkg.RequestLog
	for {
		select {
		case log := <-ls.requestsLogsService.channel:
			logs = append(logs, log)
			if len(logs) >= 5 {
				ls.flush(logs)
				logs = []*pkg.RequestLog{}
			}
		case <-ls.requestsLogsService.shutdown:
			if len(logs) > 0 {
				ls.flush(logs)
			}
			return
		}
	}

}

func (ls *LogsService) addRequestLog(log *pkg.RequestLog) error {
	ls.requestsLogsService.channel <- log
	return nil
}

func (ls *LogsService) flush(logs []*pkg.RequestLog) {
	for _, v := range logs {
		fmt.Printf("Request ID: %s\n", v.RequestID)
		fmt.Printf("Request Body: %s", v.Body)
		fmt.Println("-----")
	}
}

func (ls *LogsService) Close() error {
	close(ls.requestsLogsService.shutdown)
	ls.requestsLogsService.waitGroup.Wait()
	return nil
}
