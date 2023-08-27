package http

import (
	"docker/pkg"
	"docker/pkg/book/api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	app                  *fiber.App
	requestsResponseLogs *middleware.LogsService
	BookService          pkg.BookService
}

func New() *Server {
	s := &Server{
		app: fiber.New(),
	}

	// Logging Request ID
	s.app.Use(requestid.New())
	s.app.Use(recover.New())
	s.app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${time} ${status} - ${method} ${path}​\n",
	}))
	s.requestsResponseLogs = middleware.NewLogsService()
	s.app.Use(
		s.requestsResponseLogs.Handle(),
	)

	api := s.app.Group("/api")
	v1 := api.Group("/v1")
	s.routes(v1)

	return s
}

func (s *Server) Listen(tcpAddr string) error {
	// Create a channel to listen for an interrupt or termination signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		_ = s.app.Listen(tcpAddr)
	}()

	// Block until we receive a signal
	<-c

	// Shut down gracefully with a timeout. The timeout can be adjusted as needed
	log.Println("Gracefully shutting down...")
	s.requestsResponseLogs.Close()
	shutdownErr := s.app.Shutdown()
	if shutdownErr != nil {
		log.Fatalf("Error shutting down: %v", shutdownErr)
	}
	log.Println("Server shutdown")

	return nil
}
