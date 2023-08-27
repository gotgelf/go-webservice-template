package middleware

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/bytebufferpool"
	"net/http/httptest"
	"testing"
)

func TestHandle(t *testing.T) {

	app := fiber.New()
	app.Use(requestid.New())

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	logService := NewLogsService(Config{
		Output: buf,
	})
	app.Use(logService.Handle())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	resp, _ := app.Test(httptest.NewRequest("GET", "/test", nil))
	resp, _ = app.Test(httptest.NewRequest("GET", "/test", nil))
	resp, _ = app.Test(httptest.NewRequest("GET", "/test", nil))

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, but got %d", resp.StatusCode)
	}

	// Simulate graceful shutdown to flush any remaining logs
	logService.Close()

	// Check if logService output has our log
	if !bytes.Contains(buf.Bytes(), []byte("Request ID")) {
		t.Error("Expected logs to contain 'Request ID', but it didn't")
	}
	if !bytes.Contains(buf.Bytes(), []byte("Request Body")) {
		t.Error("Expected logs to contain 'Request Body', but it didn't")
	}
	if !bytes.Contains(buf.Bytes(), []byte("Response Status Code")) {
		t.Error("Expected logs to contain 'Response Status Code', but it didn't")
	}

	utils.AssertEqual(t, bytes.Count(buf.Bytes(), []byte("Request ID")), 3)
	utils.AssertEqual(t, bytes.Count(buf.Bytes(), []byte("Request Body")), 3)
	utils.AssertEqual(t, bytes.Count(buf.Bytes(), []byte("Response Status Code")), 3)
}
