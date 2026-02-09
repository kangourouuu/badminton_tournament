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

	// MANUAL MIGRATION: Ensure 'pool' columns exist
	// Bun's CreateSchema doesn't alter existing tables, so we do it manually here for the "v2" update.
	_, err := db.DB.Exec(`
		ALTER TABLE teams ADD COLUMN IF NOT EXISTS pool TEXT DEFAULT '';
		ALTER TABLE groups ADD COLUMN IF NOT EXISTS pool TEXT DEFAULT '';
	`)
	if err != nil {
		log.Printf("Warning: failed to alter tables: %v", err)
	}

	fmt.Println("Schema created and updated successfully")
}
