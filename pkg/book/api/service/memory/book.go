package memory

import (
	"context"
	"docker/pkg"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type BookService struct {
	DSN string
	db  map[string]*pkg.Book
	*sync.RWMutex
}

func NewBookService(DSN string) (*BookService, error) {
	return &BookService{
		DSN:     DSN,
		db:      make(map[string]*pkg.Book),
		RWMutex: &sync.RWMutex{},
	}, nil
}

func (s *BookService) CreateBook(ctx context.Context, book *pkg.Book) error {

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	bookUUID := uuid.New().String()
	book.UUID = bookUUID
	s.db[bookUUID] = book

	return nil
}

func (s *BookService) UpdateBook(ctx context.Context, book *pkg.Book) error {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	existingBook, ok := s.db[book.UUID]
	if !ok {
		return fmt.Errorf("book with UUID %s not found", book.UUID)
	}
	existingBook.Author = book.Author
	existingBook.Title = book.Title
	existingBook.Year = book.Year
	return nil
}

func (s *BookService) GetBooks(ctx context.Context) (map[string]*pkg.Book, error) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.db, nil
}

func (s *BookService) GetBook(ctx context.Context, bookUUID uuid.UUID) (*pkg.Book, error) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	existingBook, ok := s.db[bookUUID.String()]
	if !ok {
		return nil, pkg.Errorf(pkg.ENOTFOUND, "book with UUID %s not found", bookUUID.String())
	}
	return existingBook, nil
}

func (s *BookService) DeleteBook(ctx context.Context, bookUUID uuid.UUID) error {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	delete(s.db, bookUUID.String())
	return nil
}
