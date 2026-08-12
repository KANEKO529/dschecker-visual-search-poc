package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func startInferenceMock(t *testing.T, embedding []float64) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"dimension": len(embedding),
			"embedding": embedding,
		})
	}))

	t.Setenv("INFERENCE_API_BASE_URL", server.URL)

	return server
}

func TestRegisterItemEmbedding_Success(t *testing.T) {
	inference := startInferenceMock(t, []float64{0.1, 0.2, 0.3})
	defer inference.Close()

	rails := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"success": true, "embeddingId": 42})
	}))
	defer rails.Close()

	t.Setenv("RAILS_API_BASE_URL", rails.URL)
	t.Setenv("RAILS_INTERNAL_API_KEY", "test-key")

	result, err := RegisterItemEmbedding("ABC-123", strings.NewReader("fake-image"), "image.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EmbeddingID != 42 {
		t.Errorf("got embeddingId %d, want 42", result.EmbeddingID)
	}
	if result.ModelNumber != "ABC-123" {
		t.Errorf("got modelNumber %q, want ABC-123", result.ModelNumber)
	}
}

func TestRegisterItemEmbedding_InferenceUnavailable(t *testing.T) {
	// Point at a base URL nothing is listening on.
	t.Setenv("INFERENCE_API_BASE_URL", "http://localhost:0")

	_, err := RegisterItemEmbedding("ABC-123", strings.NewReader("fake-image"), "image.jpg")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var inferenceErr *InferenceError
	if !errors.As(err, &inferenceErr) {
		t.Fatalf("expected *InferenceError, got %T", err)
	}
}

func TestRegisterItemEmbedding_RailsErrorMapping(t *testing.T) {
	statuses := []int{
		http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	}

	for _, status := range statuses {
		t.Run(http.StatusText(status), func(t *testing.T) {
			inference := startInferenceMock(t, []float64{0.1, 0.2})
			defer inference.Close()

			rails := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))
			defer rails.Close()

			t.Setenv("RAILS_API_BASE_URL", rails.URL)
			t.Setenv("RAILS_INTERNAL_API_KEY", "test-key")

			_, err := RegisterItemEmbedding("ABC-123", strings.NewReader("fake-image"), "image.jpg")
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var railsErr *RailsRegistrationError
			if !errors.As(err, &railsErr) {
				t.Fatalf("expected *RailsRegistrationError, got %T", err)
			}
			if railsErr.StatusCode != status {
				t.Errorf("got status %d, want %d", railsErr.StatusCode, status)
			}
		})
	}
}
