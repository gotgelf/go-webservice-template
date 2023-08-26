package http

import "github.com/gofiber/fiber/v2"

func (s *Server) routes(router fiber.Router) {
	router.Get("/books", s.GetBooks)
	router.Get("/books/:uuid", s.GetBook)
	router.Post("/books", s.CreateBook)
	router.Put("/books/:uuid", s.UpdateBook)
	router.Delete("/books/:uuid", s.DeleteBook)
}
