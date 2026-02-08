package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/yourname/badminton-manager/backend/internal/api"
	"github.com/yourname/badminton-manager/backend/internal/db"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	r := gin.Default()
	
	// Add CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	handler := api.NewHandler(db.DB)
	handler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
