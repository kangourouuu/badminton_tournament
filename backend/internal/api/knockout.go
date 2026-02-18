package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"badminton_tournament/backend/internal/models"
	"fmt"
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

	group, err := h.EnsureKnockoutStage(c.Request.Context(), req.TournamentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created", "group_id": group.ID})
}

// EnsureKnockoutStage checks for existence and creates if missing. Returns the Group.
func (h *Handler) EnsureKnockoutStage(ctx context.Context, tournamentID uuid.UUID) (*models.Group, error) {
	// 1. Check if "KNOCKOUT" group already exists
	var existingGroup models.Group
	if err := h.DB.NewSelect().Model(&existingGroup).
		Where("tournament_id = ? AND name = ?", tournamentID, "KNOCKOUT").
		Relation("Matches").
		Scan(ctx); err == nil {
		return &existingGroup, nil
	}

	// 2. Fetch all groups and their matches to determine qualifiers
	var groups []models.Group
	err := h.DB.NewSelect().Model(&groups).
		Where("tournament_id = ?", tournamentID).
		Relation("Matches").
		Order("name ASC"). // Group A first, then Group B
		Scan(ctx)
	
	if err != nil || len(groups) < 2 {
		return nil, fmt.Errorf("Need at least 2 groups to generate knockout")
	}

	// 3. Identify Qualifiers
	type Qualifier struct {
		TeamID uuid.UUID
		Rank   int // 1 or 2
		GroupID uuid.UUID
	}
	
	qualifiers := make(map[uuid.UUID][]Qualifier)

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

	// 4. Create "KNOCKOUT" Group
	kGroup := &models.Group{
		TournamentID: tournamentID,
		Name:         "KNOCKOUT",
	}
	_, err = h.DB.NewInsert().Model(kGroup).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create knockout group: %v", err)
	}

	// 5. Create Matches (Semi-Finals & Finals)
	// Create Final & Bronze first to get IDs
	final := &models.Match{GroupID: kGroup.ID, Label: "Final"}
	_, _ = h.DB.NewInsert().Model(final).Exec(ctx)

	bronze := &models.Match{GroupID: kGroup.ID, Label: "Bronze"}
	_, _ = h.DB.NewInsert().Model(bronze).Exec(ctx)

	// Create SF1 (Cross A1 vs B2) - Placeholders only, teams filled by promotion
	sf1 := &models.Match{
		GroupID: kGroup.ID, 
		Label: "SF1", 
		NextMatchWinID: final.ID,
		NextMatchLoseID: bronze.ID,
	}
	_, _ = h.DB.NewInsert().Model(sf1).Exec(ctx)

	// Create SF2 (Cross B1 vs A2)
	sf2 := &models.Match{
		GroupID: kGroup.ID, 
		Label: "SF2", 
		NextMatchWinID: final.ID,
		NextMatchLoseID: bronze.ID,
	}
	_, _ = h.DB.NewInsert().Model(sf2).Exec(ctx)

	// Fetch fresh to return with matches
	if err := h.DB.NewSelect().Model(kGroup).Relation("Matches").WherePK().Scan(ctx); err != nil {
		return nil, fmt.Errorf("Failed to reload knockout group: %v", err)
	}
	
	return kGroup, nil
}
