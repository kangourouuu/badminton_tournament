package main

import (
	"context"
	"fmt"
	"log"

	"badminton_tournament/backend/internal/db"
	"badminton_tournament/backend/internal/models"
	
	"github.com/joho/godotenv"
)

func main() {
	// Try loading .env if present
	_ = godotenv.Load("../../.env")

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	ctx := context.Background()

	// 1. List ALL Groups
	var groups []models.Group
	if err := db.DB.NewSelect().Model(&groups).Order("name ASC").Scan(ctx); err != nil {
		log.Fatalf("Error fetching groups: %v", err)
	}

	fmt.Printf("\n=== DATABASE GROUPS LIST ===\n")
	fmt.Printf("%-36s | %-20s | %-10s\n", "Group ID", "Name", "Pool")
	fmt.Println("-----------------------------------------------------------------------")
	for _, g := range groups {
		fmt.Printf("%s | %-20s | %-10s\n", g.ID, g.Name, g.Pool)
	}
	fmt.Println("-----------------------------------------------------------------------")

	// 2. Find Knockout Group (Specific Check)
	var koGroup models.Group
	err := db.DB.NewSelect().
		Model(&koGroup).
		Where("name ILIKE ?", "knockout%").
		Relation("Matches").
		Scan(ctx)

	if err != nil {
		fmt.Printf("\n❌ CRITICAL: No 'KNOCKOUT' Group found! (Error: %v)\n", err)
		fmt.Println("You must generate the Knockout Stage (Semi-Finals) before promotion can work.")
		return
	}

	fmt.Printf("\n=== KNOCKOUT GROUP DETAILS ===\n")
	fmt.Printf("Found Group: '%s' (ID: %s)\n", koGroup.Name, koGroup.ID)
	fmt.Printf("Number of Matches: %d\n", len(koGroup.Matches))
	fmt.Println("---------------------------------------------------")
	fmt.Printf("%-36s | %-10s | %-36s | %-36s\n", "Match ID", "Label", "Team A", "Team B")
	fmt.Println("---------------------------------------------------")

	for _, m := range koGroup.Matches {
		fmt.Printf("%s | %-10s | %s | %s\n", m.ID, m.Label, m.TeamAID, m.TeamBID)
	}
	fmt.Println("---------------------------------------------------")
	
	// Check specifically for SF1 and SF2
	sf1Found := false
	sf2Found := false
	for _, m := range koGroup.Matches {
		if m.Label == "SF1" { sf1Found = true }
		if m.Label == "SF2" { sf2Found = true }
	}
	
	if sf1Found && sf2Found {
		fmt.Println("\n✅ SUCCESS: Both SF1 and SF2 labels exist.")
	} else {
		fmt.Println("\n❌ Mismatch Detected! Missing SF1 or SF2 labels.")
	}
}
