package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/service"
)

func RegisterItemEmbedding(c *gin.Context) {
	modelNumber := c.Param("modelNumber")

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

	result, err := service.RegisterItemEmbedding(modelNumber, src, file.Filename)
	if err != nil {
		status, body := mapRegisterItemEmbeddingError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":     true,
		"modelNumber": result.ModelNumber,
		"embeddingId": result.EmbeddingID,
	})
}

// mapRegisterItemEmbeddingError maps errors returned by
// service.RegisterItemEmbedding to an HTTP status and response body. Rails
// error response bodies are intentionally not forwarded to the caller; only
// the HTTP status code is used, paired with a fixed message.
func mapRegisterItemEmbeddingError(err error) (int, gin.H) {
	var inferenceErr *service.InferenceError
	if errors.As(err, &inferenceErr) {
		return http.StatusBadGateway, gin.H{"error": "failed to forward image to inference service"}
	}

	var railsErr *service.RailsRegistrationError
	if errors.As(err, &railsErr) {
		switch railsErr.StatusCode {
		case http.StatusNotFound:
			return http.StatusNotFound, gin.H{"error": "item not found for the given model number"}
		case http.StatusBadRequest:
			return http.StatusBadRequest, gin.H{"error": "invalid embedding registration request"}
		case http.StatusUnprocessableEntity:
			return http.StatusUnprocessableEntity, gin.H{"error": "invalid embedding registration request"}
		case http.StatusUnauthorized:
			return http.StatusInternalServerError, gin.H{"error": "rails internal api authentication failed"}
		default:
			return http.StatusBadGateway, gin.H{"error": "failed to register embedding in rails"}
		}
	}

	return http.StatusBadGateway, gin.H{"error": "failed to register embedding"}
}
