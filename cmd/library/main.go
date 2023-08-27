package main

import (
	"docker/pkg/book/api/http"
	"docker/pkg/book/api/service/memory"
	"fmt"
	"log"
)

func main() {
	server := http.New()

	bookService, err := memory.NewBookService("memory")
	if err != nil {
		fmt.Errorf("failed to create book service: %v", err)
	}

	server.BookService = bookService

	log.Fatal(server.Listen(":8080"))
}
