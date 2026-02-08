package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"badminton_tournament/backend/internal/api"
	"badminton_tournament/backend/internal/auth"
	"badminton_tournament/backend/internal/db"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	
	// Init Auth
	auth.Init()

	r := gin.Default()
	
	// Add CORS
	config := cors.DefaultConfig()
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		config.AllowOrigins = []string{allowedOrigins}
	} else {
		config.AllowAllOrigins = true
	}
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
