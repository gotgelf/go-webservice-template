package pkg

import (
	"context"
	"github.com/google/uuid"
)

type Book struct {
	UUID   string `json:"uuid,omitempty"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year,omitempty"`
}

type BookService interface {
	CreateBook(ctx context.Context, book *Book) error
	UpdateBook(ctx context.Context, book *Book) error
	GetBooks(ctx context.Context) (map[string]*Book, error)
	GetBook(ctx context.Context, bookUUID uuid.UUID) (*Book, error)
	DeleteBook(ctx context.Context, bookUUID uuid.UUID) error
}
