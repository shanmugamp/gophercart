package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	HealthHandler(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	if response.Body.String() != `{"status":"ok"}` {
		t.Errorf("Expected body `{'status':'ok'}`, got `%s`", response.Body.String())
	}
}
