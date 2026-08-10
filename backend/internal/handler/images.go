package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/client"
)

func UploadImage(c *gin.Context) {
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

	if err := client.ForwardImage(src, file.Filename); err != nil {
		c.JSON(502, gin.H{"error": "failed to forward image to inference service"})
		return
	}

	c.JSON(200, gin.H{
		"filename":    file.Filename,
		"size":        file.Size,
		"contentType": file.Header.Get("Content-Type"),
	})
}
