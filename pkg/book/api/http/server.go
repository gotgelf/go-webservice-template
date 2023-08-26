package http

import (
	"docker/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	app         *fiber.App
	BookService pkg.BookService
}

func New() *Server {
	s := &Server{
		app: fiber.New(),
	}

	s.app.Use(recover.New())
	s.app.Use(logger.New())

	api := s.app.Group("/api")
	v1 := api.Group("/v1")
	s.routes(v1)

	return s
}

func (s *Server) Listen(tcpAddr string) error {
	return s.app.Listen(tcpAddr)
}
