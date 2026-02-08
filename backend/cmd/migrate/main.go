package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourname/badminton-manager/backend/internal/db"
	"github.com/yourname/badminton-manager/backend/internal/models"
)

func CreateSchema(ctx context.Context) error {
	modelsToRegister := []interface{}{
		(*models.Participant)(nil),
		(*models.Team)(nil),
		(*models.Match)(nil),
	}

	for _, model := range modelsToRegister {
		_, err := db.DB.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return fmt.Errorf("creating table: %w", err)
		}
	}
	return nil
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := CreateSchema(context.Background()); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	fmt.Println("Schema created successfully")
}
