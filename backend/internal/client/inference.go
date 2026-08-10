package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const inferenceBaseURL = "http://localhost:3002"

var httpClient = &http.Client{Timeout: 10 * time.Second}

// ForwardImage sends the image to the inference service's image-receiving
// endpoint as multipart/form-data. It returns an error if the request could
// not be sent or the inference service responded with a non-2xx status.
func ForwardImage(file io.Reader, filename string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("copy image content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, inferenceBaseURL+"/internal/images", body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("inference service returned status %d", resp.StatusCode)
	}

	return nil
}
