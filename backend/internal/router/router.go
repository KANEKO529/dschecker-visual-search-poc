package router

import (
	"github.com/gin-gonic/gin"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/handler"
)

func New() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.Health)
	r.POST("/api/poc/images", handler.UploadImage)

	return r
}
