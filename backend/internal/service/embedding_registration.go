package service

import (
	"errors"
	"fmt"
	"io"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/client"
)

// InferenceError indicates the Python inference service could not be
// reached, or returned an error while generating the embedding.
type InferenceError struct {
	Err error
}

func (e *InferenceError) Error() string {
	return fmt.Sprintf("generate embedding: %v", e.Err)
}

func (e *InferenceError) Unwrap() error {
	return e.Err
}

// RailsRegistrationError indicates the Rails internal API returned a
// non-success response while registering the embedding. StatusCode is the
// HTTP status code Rails responded with, or 0 if Rails could not be reached
// at all.
type RailsRegistrationError struct {
	StatusCode int
	Err        error
}

func (e *RailsRegistrationError) Error() string {
	return fmt.Sprintf("register embedding in rails: %v", e.Err)
}

func (e *RailsRegistrationError) Unwrap() error {
	return e.Err
}

// RegisterItemEmbeddingResult is the outcome of successfully registering an
// item's embedding in Rails.
type RegisterItemEmbeddingResult struct {
	ModelNumber string
	EmbeddingID int
}

// RegisterItemEmbedding generates an embedding for the given image via the
// Python inference service, then registers it against the item identified
// by modelNumber in Rails.
func RegisterItemEmbedding(modelNumber string, image io.Reader, filename string) (*RegisterItemEmbeddingResult, error) {
	embedding, err := client.GenerateEmbedding(image, filename)
	if err != nil {
		return nil, &InferenceError{Err: err}
	}

	resp, err := client.RegisterEmbedding(modelNumber, embedding.Embedding)
	if err != nil {
		var railsErr *client.RailsAPIError
		statusCode := 0
		if errors.As(err, &railsErr) {
			statusCode = railsErr.StatusCode
		}
		return nil, &RailsRegistrationError{StatusCode: statusCode, Err: err}
	}

	return &RegisterItemEmbeddingResult{
		ModelNumber: modelNumber,
		EmbeddingID: resp.EmbeddingID,
	}, nil
}
