package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/service"
)

func TestSearchItemsByImage_MissingImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/poc/items/search-by-image", SearchItemsByImage)

	req := httptest.NewRequest(http.MethodPost, "/api/poc/items/search-by-image", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestMapSearchItemsByImageError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"inference error", &service.InferenceError{Err: errors.New("boom")}, http.StatusBadGateway},
		{"rails bad request", &service.RailsSearchError{StatusCode: http.StatusBadRequest, Err: errors.New("x")}, http.StatusBadGateway},
		{"rails unauthorized", &service.RailsSearchError{StatusCode: http.StatusUnauthorized, Err: errors.New("x")}, http.StatusInternalServerError},
		{"rails internal server error", &service.RailsSearchError{StatusCode: http.StatusInternalServerError, Err: errors.New("x")}, http.StatusBadGateway},
		{"unknown error", errors.New("unknown"), http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, _ := mapSearchItemsByImageError(tt.err)
			if status != tt.wantStatus {
				t.Errorf("got status %d, want %d", status, tt.wantStatus)
			}
		})
	}
}
