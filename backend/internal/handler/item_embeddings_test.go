package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/service"
)

func TestRegisterItemEmbedding_MissingImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/poc/items/:modelNumber/embeddings", RegisterItemEmbedding)

	req := httptest.NewRequest(http.MethodPost, "/api/poc/items/ABC-123/embeddings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestMapRegisterItemEmbeddingError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"inference error", &service.InferenceError{Err: errors.New("boom")}, http.StatusBadGateway},
		{"rails not found", &service.RailsRegistrationError{StatusCode: http.StatusNotFound, Err: errors.New("x")}, http.StatusNotFound},
		{"rails bad request", &service.RailsRegistrationError{StatusCode: http.StatusBadRequest, Err: errors.New("x")}, http.StatusBadRequest},
		{"rails unprocessable entity", &service.RailsRegistrationError{StatusCode: http.StatusUnprocessableEntity, Err: errors.New("x")}, http.StatusUnprocessableEntity},
		{"rails internal server error", &service.RailsRegistrationError{StatusCode: http.StatusInternalServerError, Err: errors.New("x")}, http.StatusBadGateway},
		{"rails unauthorized", &service.RailsRegistrationError{StatusCode: http.StatusUnauthorized, Err: errors.New("x")}, http.StatusInternalServerError},
		{"unknown error", errors.New("unknown"), http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, _ := mapRegisterItemEmbeddingError(tt.err)
			if status != tt.wantStatus {
				t.Errorf("got status %d, want %d", status, tt.wantStatus)
			}
		})
	}
}
