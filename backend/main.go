package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found, relying on environment variables: %v", err)
	}

	r := router.New()
	r.Run("localhost:3001")
}
