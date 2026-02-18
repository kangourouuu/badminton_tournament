package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for local testing if needed, though strictly we expect env var
		dsn = "postgres://user:password@localhost:5432/badminton?sslmode=disable"
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	ctx := context.Background()

	// Add partner_request column
	_, err := db.ExecContext(ctx, "ALTER TABLE participants ADD COLUMN IF NOT EXISTS partner_request TEXT;")
	if err != nil {
		log.Printf("Error adding column: %v", err)
	} else {
		fmt.Println("Successfully added 'partner_request' column to participants table.")
	}
}
