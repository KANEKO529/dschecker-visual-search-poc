package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// RailsAPIError indicates the Rails internal API responded with a non-2xx
// HTTP status.
type RailsAPIError struct {
	StatusCode int
}

func (e *RailsAPIError) Error() string {
	return fmt.Sprintf("rails internal api returned status %d", e.StatusCode)
}

// RailsEmbeddingResponse is the response body returned by the Rails internal
// embedding registration API on success.
type RailsEmbeddingResponse struct {
	Success     bool `json:"success"`
	EmbeddingID int  `json:"embeddingId"`
}

type railsEmbeddingRequest struct {
	Embedding []float64 `json:"embedding"`
}

// RailsSearchResult is a single item returned by the Rails internal
// embedding-based search API, along with its similarity to the query
// embedding.
type RailsSearchResult struct {
	ModelNumber string  `json:"modelNumber"`
	Similarity  float64 `json:"similarity"`
}

// RailsSearchResponse is the response body returned by the Rails internal
// embedding-based search API on success.
type RailsSearchResponse struct {
	Results []RailsSearchResult `json:"results"`
}

type railsSearchRequest struct {
	Embedding []float64 `json:"embedding"`
}

// RegisterEmbedding sends the generated embedding to the Rails internal API
// to be registered against the item identified by modelNumber. It returns a
// *RailsAPIError if Rails responds with a non-2xx status.
func RegisterEmbedding(modelNumber string, embedding []float64) (*RailsEmbeddingResponse, error) {
	baseURL := os.Getenv("RAILS_API_BASE_URL")
	apiKey := os.Getenv("RAILS_INTERNAL_API_KEY")

	reqBody, err := json.Marshal(railsEmbeddingRequest{Embedding: embedding})
	if err != nil {
		return nil, fmt.Errorf("encode request body: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/internal/items/%s/embeddings", baseURL, modelNumber)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-API-Key", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &RailsAPIError{StatusCode: resp.StatusCode}
	}

	var embeddingResp RailsEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &embeddingResp, nil
}

// SearchByEmbedding sends the generated embedding to the Rails internal API
// to search for items with similar embeddings. It returns a *RailsAPIError
// if Rails responds with a non-2xx status.
func SearchByEmbedding(embedding []float64) (*RailsSearchResponse, error) {
	baseURL := os.Getenv("RAILS_API_BASE_URL")
	apiKey := os.Getenv("RAILS_INTERNAL_API_KEY")

	reqBody, err := json.Marshal(railsSearchRequest{Embedding: embedding})
	if err != nil {
		return nil, fmt.Errorf("encode request body: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/internal/items/search_by_embedding", baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-API-Key", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &RailsAPIError{StatusCode: resp.StatusCode}
	}

	var searchResp RailsSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &searchResp, nil
}
