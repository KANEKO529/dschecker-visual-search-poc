package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/service"
)

func SearchItemsByImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image"})
		return
	}
	defer src.Close()

	result, err := service.SearchItemsByImage(src, file.Filename)
	if err != nil {
		status, body := mapSearchItemsByImageError(err)
		c.JSON(status, body)
		return
	}

	results := make([]gin.H, len(result.Results))
	for i, r := range result.Results {
		results[i] = gin.H{
			"modelNumber": r.ModelNumber,
			"similarity":  r.Similarity,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}

// mapSearchItemsByImageError maps errors returned by
// service.SearchItemsByImage to an HTTP status and response body. Rails
// error response bodies are intentionally not forwarded to the caller; only
// the HTTP status code is used, paired with a fixed message.
func mapSearchItemsByImageError(err error) (int, gin.H) {
	var inferenceErr *service.InferenceError
	if errors.As(err, &inferenceErr) {
		return http.StatusBadGateway, gin.H{"error": "failed to forward image to inference service"}
	}

	var railsErr *service.RailsSearchError
	if errors.As(err, &railsErr) {
		switch railsErr.StatusCode {
		case http.StatusUnauthorized:
			return http.StatusInternalServerError, gin.H{"error": "rails internal api authentication failed"}
		default:
			return http.StatusBadGateway, gin.H{"error": "failed to search items in rails"}
		}
	}

	return http.StatusBadGateway, gin.H{"error": "failed to search items by image"}
}
