package http

import (
	"bytes"
	"docker/pkg"
	"docker/pkg/book/api/service/memory"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBooks(t *testing.T) {
	s := New()
	bookService, _ := memory.NewBookService("memory")
	s.BookService = bookService

	req := httptest.NewRequest("GET", "/api/v1/books", nil)

	resp, err := s.app.Test(req)
	if err != nil {
		t.Fatalf("failed to get response: %v", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		t.Fatalf("expected code: %d, got: %d", http.StatusOK, code)
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
}

func TestCreateBook(t *testing.T) {
	s := New()
	bookService, _ := memory.NewBookService("memory")
	s.BookService = bookService

	book := `{	
	"title": "title1",
	"author": "author1",
	"year": 2012
}`
	req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewBufferString(book))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	if err != nil {
		t.Fatalf("failed to get response: %v", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		t.Fatalf("expected code: %d, got: %d", http.StatusOK, code)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	respBook := &pkg.Book{}

	if err := json.Unmarshal(body, respBook); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if respBook.Title != "title1" {
		t.Fatalf("expected name: %s, got: %s", "title1", respBook.Title)
	}
}

func TestGetBook(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		s := New()

		bookService, _ := memory.NewBookService("memory")
		s.BookService = bookService

		book := `{	
	"title": "title1",
	"author": "author1",
	"year": 2012
}`
		req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewBufferString(book))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		respBook := &pkg.Book{}

		if err := json.Unmarshal(body, respBook); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		UUID := respBook.UUID
		req = httptest.NewRequest("GET", "/api/v1/books/"+UUID, nil)
		resp, err = s.app.Test(req)

		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			t.Fatalf("expected code: %d, got: %d", http.StatusOK, code)
		}

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		respBook = &pkg.Book{}
		if err := json.Unmarshal(body, respBook); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if respBook.UUID != UUID {
			t.Fatalf("expected UUID: %s, got: %s", UUID, respBook.UUID)
		}
	})

	t.Run("404", func(t *testing.T) {
		s := New()

		bookService, _ := memory.NewBookService("memory")
		s.BookService = bookService

		book := `{	
	"title": "title1",
	"author": "author1",
	"year": 2012
}`
		req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewBufferString(book))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		respBook := &pkg.Book{}

		if err := json.Unmarshal(body, respBook); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		notExistingUUID := respBook.UUID + "not-existing"
		req = httptest.NewRequest("GET", "/api/v1/books/"+notExistingUUID, nil)
		resp, err = s.app.Test(req)

		if err != nil {
			t.Fatalf("failed to get response: %v", err)
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusInternalServerError {
			t.Fatalf("expected code: %d, got: %d", http.StatusInternalServerError, code)
		}
	})
}
