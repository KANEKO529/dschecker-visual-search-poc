package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/client"
)

func GenerateEmbedding(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read image"})
		return
	}
	defer src.Close()

	embedding, err := client.GenerateEmbedding(src, file.Filename)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to forward image to inference service"})
		return
	}

	c.JSON(200, gin.H{
		"dimension": embedding.Dimension,
		"embedding": embedding.Embedding,
	})
}
