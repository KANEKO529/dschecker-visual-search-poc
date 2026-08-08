package main

import "github.com/KANEKO529/dschecker-visual-search-poc/backend/internal/router"

func main() {
	r := router.New()
	r.Run("localhost:3001")
}
