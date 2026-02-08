package main

import (
	"context"
	"fmt"
	"log"

	"badminton_tournament/backend/internal/db"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := db.CreateSchema(context.Background()); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	fmt.Println("Schema created successfully")
}
