package api

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"badminton_tournament/backend/internal/models"
)

// GenerateKnockoutRequest
type GenerateKnockoutRequest struct {
	TournamentID uuid.UUID `json:"tournament_id"`
}

func (h *Handler) GenerateKnockout(c *gin.Context) {
	var req GenerateKnockoutRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. Check if "KNOCKOUT" group already exists
	count, _ := h.DB.NewSelect().Model((*models.Group)(nil)).
		Where("tournament_id = ? AND name = ?", req.TournamentID, "KNOCKOUT").
		Count(ctx)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Knockout stage already exists"})
		return
	}

	// 2. Fetch all groups and their matches to determine qualifiers
	// Assumption: Group A and Group B are the only groups, or we take top 2 from any 2 groups found.
	// For simplicity in this 1-day build, let's assume standard 2-group format (Mesoneer logic).
	var groups []models.Group
	err := h.DB.NewSelect().Model(&groups).
		Where("tournament_id = ?", req.TournamentID).
		Relation("Matches").
		Order("name ASC"). // Group A first, then Group B
		Scan(ctx)
	
	if err != nil || len(groups) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Need at least 2 groups to generate knockout"})
		return
	}

	// 3. Identify Qualifiers
	// Logic:
	// - Winner of M3 (Winners Match) -> 1st Place
	// - Winner of M5 (Decider Match) -> 2nd Place
	
	type Qualifier struct {
		TeamID uuid.UUID
		Rank   int // 1 or 2
		GroupID uuid.UUID
	}
	
	qualifiers := make(map[uuid.UUID][]Qualifier) // GroupID -> []Qualifier

	for _, g := range groups {
		var first, second uuid.UUID
		for _, m := range g.Matches {
			if m.Label == "Winners" { // M3
				first = m.WinnerID
			}
			if m.Label == "Decider" { // M5
				second = m.WinnerID
			}
		}
		
		if first != uuid.Nil {
			qualifiers[g.ID] = append(qualifiers[g.ID], Qualifier{TeamID: first, Rank: 1, GroupID: g.ID})
		}
		if second != uuid.Nil {
			qualifiers[g.ID] = append(qualifiers[g.ID], Qualifier{TeamID: second, Rank: 2, GroupID: g.ID})
		}
	}

	// Validate we have enough qualifiers
	// We need A1, A2, B1, B2
	groupA := groups[0]
	groupB := groups[1]
	
	if len(qualifiers[groupA.ID]) < 2 || len(qualifiers[groupB.ID]) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Matches not finished. Determine 1st and 2nd place for all groups."})
		return
	}

	var a1, a2, b1, b2 uuid.UUID
	
	// Map properly
	for _, q := range qualifiers[groupA.ID] {
		if q.Rank == 1 { a1 = q.TeamID } else { a2 = q.TeamID }
	}
	for _, q := range qualifiers[groupB.ID] {
		if q.Rank == 1 { b1 = q.TeamID } else { b2 = q.TeamID }
	}

	// 4. Create "KNOCKOUT" Group
	kGroup := &models.Group{
		TournamentID: req.TournamentID,
		Name:         "KNOCKOUT",
	}
	_, err = h.DB.NewInsert().Model(kGroup).Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create knockout group"})
		return
	}

	// 5. Create Matches (Semi-Finals & Finals)
	// Structure:
	// SF1: A1 vs B2 -> Win: Final, Lose: Bronze
	// SF2: B1 vs A2 -> Win: Final, Lose: Bronze
	// Bronze: Loser SF1 vs Loser SF2
	// Final: Winner SF1 vs Winner SF2

	// Create Final & Bronze first to get IDs
	final := &models.Match{GroupID: kGroup.ID, Label: "Final"}
	_, _ = h.DB.NewInsert().Model(final).Exec(ctx)

	bronze := &models.Match{GroupID: kGroup.ID, Label: "Bronze"}
	_, _ = h.DB.NewInsert().Model(bronze).Exec(ctx)

	// Create SF1 (Cross A1 vs B2)
	sf1 := &models.Match{
		GroupID: kGroup.ID, 
		Label: "SF1", 
		TeamAID: a1, 
		TeamBID: b2,
		NextMatchWinID: final.ID,
		NextMatchLoseID: bronze.ID,
	}
	_, _ = h.DB.NewInsert().Model(sf1).Exec(ctx)

	// Create SF2 (Cross B1 vs A2)
	sf2 := &models.Match{
		GroupID: kGroup.ID, 
		Label: "SF2", 
		TeamAID: b1, 
		TeamBID: a2,
		NextMatchWinID: final.ID,
		NextMatchLoseID: bronze.ID,
	}
	_, _ = h.DB.NewInsert().Model(sf2).Exec(ctx)

	c.JSON(http.StatusOK, gin.H{"status": "created", "group_id": kGroup.ID})
}
