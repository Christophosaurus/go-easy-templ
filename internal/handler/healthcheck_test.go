package handler_test

import (
	"github.com/go-easy-templ/internal/config"
	"github.com/go-easy-templ/internal/handler"
	"github.com/go-easy-templ/internal/logger"
	"github.com/go-easy-templ/internal/repository"
	"github.com/go-easy-templ/internal/service"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {

	// Create mocks
	mockConfig := &config.Config{
		Server: config.Server{
			Env: "test",
		},
	}

	mockLogger := logger.NewSlogger(logger.LevelDebug)

	mockDummyRepo := repository.NewDummyRepository(nil)
	mockRepos := repository.NewRepositories(mockDummyRepo)
	mockService := service.NewDummy(nil, mockLogger, mockRepos)
	mockServices := service.NewServices(mockService)
	healthcheckHandler := handler.NewHealthcheckHandler(mockLogger, mockConfig, mockServices)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	healthcheckHandler.Healthcheck(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Read response body
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	// Check the response body is what we expect.
	actual := strings.TrimRight(string(b), "\n")
	expected := `{"status":"available","system_info":{"environment":"test","version":"1.0.0"}}`

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
