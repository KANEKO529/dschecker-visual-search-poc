package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterEmbedding_Success(t *testing.T) {
	var gotPath, gotAPIKey, gotContentType string
	var gotBody railsEmbeddingRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAPIKey = r.Header.Get("X-Internal-API-Key")
		gotContentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(RailsEmbeddingResponse{Success: true, EmbeddingID: 7})
	}))
	defer server.Close()

	t.Setenv("RAILS_API_BASE_URL", server.URL)
	t.Setenv("RAILS_INTERNAL_API_KEY", "secret-key")

	resp, err := RegisterEmbedding("ABC-123", []float64{0.1, 0.2, 0.3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.EmbeddingID != 7 {
		t.Errorf("got embeddingId %d, want 7", resp.EmbeddingID)
	}
	if gotPath != "/api/v1/internal/items/ABC-123/embeddings" {
		t.Errorf("got path %q, want /api/v1/internal/items/ABC-123/embeddings", gotPath)
	}
	if gotAPIKey != "secret-key" {
		t.Errorf("got X-Internal-API-Key %q, want secret-key", gotAPIKey)
	}
	if gotContentType != "application/json" {
		t.Errorf("got Content-Type %q, want application/json", gotContentType)
	}
	if len(gotBody.Embedding) != 3 {
		t.Errorf("got embedding length %d, want 3", len(gotBody.Embedding))
	}
}

func TestRegisterEmbedding_ErrorStatuses(t *testing.T) {
	statuses := []int{
		http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	}

	for _, status := range statuses {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))
			defer server.Close()

			t.Setenv("RAILS_API_BASE_URL", server.URL)
			t.Setenv("RAILS_INTERNAL_API_KEY", "secret-key")

			_, err := RegisterEmbedding("ABC-123", []float64{0.1})
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var railsErr *RailsAPIError
			if !errors.As(err, &railsErr) {
				t.Fatalf("expected *RailsAPIError, got %T", err)
			}
			if railsErr.StatusCode != status {
				t.Errorf("got status %d, want %d", railsErr.StatusCode, status)
			}
		})
	}
}
