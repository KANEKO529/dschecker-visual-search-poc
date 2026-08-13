package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchItemsByImage_Success(t *testing.T) {
	inference := startInferenceMock(t, []float64{0.1, 0.2, 0.3})
	defer inference.Close()

	rails := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{"modelNumber": "ABC-123", "similarity": 0.95},
			},
		})
	}))
	defer rails.Close()

	t.Setenv("RAILS_API_BASE_URL", rails.URL)
	t.Setenv("RAILS_INTERNAL_API_KEY", "test-key")

	result, err := SearchItemsByImage(strings.NewReader("fake-image"), "image.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Results) != 1 {
		t.Fatalf("got %d results, want 1", len(result.Results))
	}
	if result.Results[0].ModelNumber != "ABC-123" || result.Results[0].Similarity != 0.95 {
		t.Errorf("got result[0] %+v, want {ABC-123 0.95}", result.Results[0])
	}
}

func TestSearchItemsByImage_EmptyResults(t *testing.T) {
	inference := startInferenceMock(t, []float64{0.1, 0.2, 0.3})
	defer inference.Close()

	rails := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{}})
	}))
	defer rails.Close()

	t.Setenv("RAILS_API_BASE_URL", rails.URL)
	t.Setenv("RAILS_INTERNAL_API_KEY", "test-key")

	result, err := SearchItemsByImage(strings.NewReader("fake-image"), "image.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Results) != 0 {
		t.Errorf("got %d results, want 0", len(result.Results))
	}
}

func TestSearchItemsByImage_InferenceUnavailable(t *testing.T) {
	// Point at a base URL nothing is listening on.
	t.Setenv("INFERENCE_API_BASE_URL", "http://localhost:0")

	_, err := SearchItemsByImage(strings.NewReader("fake-image"), "image.jpg")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var inferenceErr *InferenceError
	if !errors.As(err, &inferenceErr) {
		t.Fatalf("expected *InferenceError, got %T", err)
	}
}

func TestSearchItemsByImage_RailsErrorMapping(t *testing.T) {
	statuses := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
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

			_, err := SearchItemsByImage(strings.NewReader("fake-image"), "image.jpg")
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var railsErr *RailsSearchError
			if !errors.As(err, &railsErr) {
				t.Fatalf("expected *RailsSearchError, got %T", err)
			}
			if railsErr.StatusCode != status {
				t.Errorf("got status %d, want %d", railsErr.StatusCode, status)
			}
		})
	}
}
