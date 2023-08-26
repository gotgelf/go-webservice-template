package http

import (
	"context"
	"docker/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

func (s *Server) CreateBook(c *fiber.Ctx) error {
	book := new(pkg.Book)
	if err := c.BodyParser(book); err != nil {
		return err
	}

	err := s.BookService.CreateBook(context.TODO(), book)
	if err != nil {
		return err
	}

	return c.JSON(book)
}

func (s *Server) UpdateBook(c *fiber.Ctx) error {
	bookUUID, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return err
	}

	book := new(pkg.Book)
	if err := c.BodyParser(book); err != nil {
		return err
	}
	book.UUID = bookUUID.String()

	err = s.BookService.UpdateBook(context.TODO(), book)
	if err != nil {
		return err
	}

	return c.JSON(book)
}

func (s *Server) GetBooks(c *fiber.Ctx) error {
	books, err := s.BookService.GetBooks(context.TODO())
	if err != nil {
		return err
	}

	return c.JSON(books)
}

func (s *Server) GetBook(c *fiber.Ctx) error {
	bookUUID, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return err
	}

	book, err := s.BookService.GetBook(context.TODO(), bookUUID)
	if err != nil {
		return err
	}

	return c.JSON(book)
}

func (s *Server) DeleteBook(c *fiber.Ctx) error {
	bookUUID, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return err
	}

	err = s.BookService.DeleteBook(context.TODO(), bookUUID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK)
}
