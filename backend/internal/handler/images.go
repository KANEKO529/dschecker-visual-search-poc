package handler

import "github.com/gin-gonic/gin"

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image is required"})
		return
	}

	c.JSON(200, gin.H{
		"filename":    file.Filename,
		"size":        file.Size,
		"contentType": file.Header.Get("Content-Type"),
	})
}
