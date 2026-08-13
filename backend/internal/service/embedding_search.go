package service

import (
	"errors"
	"fmt"
	"io"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/client"
)

// RailsSearchError indicates the Rails internal API returned a non-success
// response while searching for items by embedding. StatusCode is the HTTP
// status code Rails responded with, or 0 if Rails could not be reached at
// all.
type RailsSearchError struct {
	StatusCode int
	Err        error
}

func (e *RailsSearchError) Error() string {
	return fmt.Sprintf("search items by embedding in rails: %v", e.Err)
}

func (e *RailsSearchError) Unwrap() error {
	return e.Err
}

// SearchResult is a single item returned by an image search, along with its
// similarity to the query image.
type SearchResult struct {
	ModelNumber string
	Similarity  float64
}

// SearchItemsByImageResult is the outcome of successfully searching for
// items by image.
type SearchItemsByImageResult struct {
	Results []SearchResult
}

// SearchItemsByImage generates an embedding for the given image via the
// Python inference service, then searches for items with similar embeddings
// in Rails.
func SearchItemsByImage(image io.Reader, filename string) (*SearchItemsByImageResult, error) {
	embedding, err := client.GenerateEmbedding(image, filename)
	if err != nil {
		return nil, &InferenceError{Err: err}
	}

	resp, err := client.SearchByEmbedding(embedding.Embedding)
	if err != nil {
		var railsErr *client.RailsAPIError
		statusCode := 0
		if errors.As(err, &railsErr) {
			statusCode = railsErr.StatusCode
		}
		return nil, &RailsSearchError{StatusCode: statusCode, Err: err}
	}

	results := make([]SearchResult, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = SearchResult{ModelNumber: r.ModelNumber, Similarity: r.Similarity}
	}

	return &SearchItemsByImageResult{Results: results}, nil
}
