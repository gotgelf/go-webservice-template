package middleware

import (
	"docker/pkg"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sync"
)

type LogsService struct {
	channel   chan *pkg.Log
	shutdown  chan bool
	waitGroup sync.WaitGroup
}

func NewLogsService() *LogsService {
	service := &LogsService{
		channel:  make(chan *pkg.Log, 5),
		shutdown: make(chan bool),
	}
	service.waitGroup.Add(1)

	go service.listen()

	return service
}

func (ls *LogsService) Handle() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		body := string(c.Request().Body())
		requestID := c.Locals("requestid").(string)
		_ = ls.addLog(&pkg.Log{
			RequestID:          requestID,
			Body:               body,
			ResponseStatusCode: c.Response().StatusCode(),
		})

		return c.Next()
	}
}

func (ls *LogsService) listen() {
	defer ls.waitGroup.Done()
	var requestsLogs []*pkg.Log
	for {
		select {
		case log := <-ls.channel:
			requestsLogs = append(requestsLogs, log)
			if len(requestsLogs) >= 5 {
				ls.flushRequestsLogs(requestsLogs)
				requestsLogs = []*pkg.Log{}
			}
		case <-ls.shutdown:
			if len(requestsLogs) > 0 {
				ls.flushRequestsLogs(requestsLogs)
			}
			return
		}
	}
}

func (ls *LogsService) addLog(log *pkg.Log) error {
	ls.channel <- log
	return nil
}

func (ls *LogsService) flushRequestsLogs(logs []*pkg.Log) {
	for _, v := range logs {
		fmt.Printf("Request ID: %s\n", v.RequestID)
		fmt.Printf("Request Body: %s", v.Body)
		fmt.Printf("Response Status Code: %d", v.ResponseStatusCode)
		fmt.Println("-----")
	}
}

func (ls *LogsService) Close() error {
	close(ls.shutdown)
	ls.waitGroup.Wait()
	return nil
}
