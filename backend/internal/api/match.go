package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"badminton_tournament/backend/internal/models"
)

type UpdateMatchRequest struct {
	WinnerID uuid.UUID `json:"winner_id"`
	Score    string    `json:"score"`
	VideoURL string    `json:"video_url"`
}

func (h *Handler) UpdateMatch(c *gin.Context) {
	id := c.Param("id")
	var req UpdateMatchRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. Get current match
	var match models.Match
	err := h.DB.NewSelect().Model(&match).Where("id = ?", id).Scan(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// 2. Update current match
	match.WinnerID = req.WinnerID
	match.Score = req.Score
	match.VideoURL = req.VideoURL
	
	_, err = h.DB.NewUpdate().Model(&match).Column("winner_id", "score", "video_url").WherePK().Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Auto-Propagation (The Magic)
	// Only proceed if we have a winner
	if req.WinnerID != uuid.Nil {
		// Identify Loser
		var loserID uuid.UUID
		if match.TeamAID == req.WinnerID {
			loserID = match.TeamBID
		} else {
			loserID = match.TeamAID
		}

		// Propagate Winner
		if match.NextMatchWinID != uuid.Nil {
			h.propagateToMatch(ctx, match.NextMatchWinID, req.WinnerID, match.Label, "win")
		}

		// Propagate Loser
		if match.NextMatchLoseID != uuid.Nil && loserID != uuid.Nil {
			h.propagateToMatch(ctx, match.NextMatchLoseID, loserID, match.Label, "lose")
		}

		// 4. Auto-Flow: Check for Qualification (GSL -> Knockout)
		// Winners Match (M3) Winner -> Rank 1
		if match.Label == "Winners" {
			h.promoteToKnockout(ctx, match.GroupID, 1, req.WinnerID)
		}
		// Decider Match (M5) Winner -> Rank 2
		if match.Label == "Decider" {
			h.promoteToKnockout(ctx, match.GroupID, 2, req.WinnerID)
		}
	}

	c.JSON(http.StatusOK, match)
}

func (h *Handler) propagateToMatch(ctx context.Context, targetID, teamID uuid.UUID, sourceLabel, outcome string) error {
	// Find target match
	var target models.Match
	err := h.DB.NewSelect().Model(&target).Where("id = ?", targetID).Scan(ctx)
	if err != nil {
		return err
	}

	col := ""

	// Deterministic GSL Logic
	// M1 -> Win -> M3 (Slot A)
	// M2 -> Win -> M3 (Slot B)
	// M1 -> Lose -> M4 (Slot A)
	// M2 -> Lose -> M4 (Slot B)
	// M3 -> Lose -> M5 (Slot A)
	// M4 -> Win -> M5 (Slot B)

	if target.Label == "Winners" { // M3
		if sourceLabel == "M1" {
			col = "team_a_id"
			target.TeamAID = teamID
		} else if sourceLabel == "M2" {
			col = "team_b_id"
			target.TeamBID = teamID
		}
	} else if target.Label == "Losers" { // M4
		if sourceLabel == "M1" {
			col = "team_a_id"
			target.TeamAID = teamID
		} else if sourceLabel == "M2" {
			col = "team_b_id"
			target.TeamBID = teamID
		}
	} else if target.Label == "Decider" { // M5
		if sourceLabel == "Winners" { // From M3 Loser
			col = "team_a_id"
			target.TeamAID = teamID
		} else if sourceLabel == "Losers" { // From M4 Winner
			col = "team_b_id"
			target.TeamBID = teamID
		}
	}

	// Fallback to "First Empty" if label logic doesn't match (e.g. custom bracket)
	if col == "" {
		if target.TeamAID == uuid.Nil {
			target.TeamAID = teamID
			col = "team_a_id"
		} else if target.TeamBID == uuid.Nil {
			target.TeamBID = teamID
			col = "team_b_id"
		}
	}

// ... (propagateToMatch existing code)
	if col != "" {
		_, err = h.DB.NewUpdate().Model(&target).Column(col).WherePK().Exec(ctx)
		return err
	}
	return nil
}

// Auto-Flow: Promote Group Winners to Knockout Bracket
func (h *Handler) promoteToKnockout(ctx context.Context, groupID uuid.UUID, rank int, teamID uuid.UUID) error {
	// 1. Get Source Group Name (to determine A vs B)
	var group models.Group
	if err := h.DB.NewSelect().Model(&group).Where("id = ?", groupID).Scan(ctx); err != nil {
		return err
	}

	// 2. Get Knockout Group & Matches
	var koGroup models.Group
	if err := h.DB.NewSelect().Model(&koGroup).Where("name = ?", "KNOCKOUT").Relation("Matches").Scan(ctx); err != nil {
		return err // Knockout bracket might not exist yet
	}

	// 3. Find Target Slot
	// Logic:
	// Group A (Mesoneer?) -> "Group A"
	// Group B (Lab?)      -> "Group B"
	// Mapping:
	// A1 -> SF1 (Home)
	// B2 -> SF1 (Away)
	// B1 -> SF2 (Home)
	// A2 -> SF2 (Away)

	var targetLabel string
	var targetCol string

	// Heuristic: Check if name contains "A" or "Mesoneer" vs "B" or "Lab"
	// NOTE: This relies on naming conventions "Group A" / "Group B" or Pools.
	// Let's use Pool if available, or Name.
	isGroupA := group.Pool == "Mesoneer" || group.Name == "Group A"
	// isGroupB := group.Pool == "Lab" || group.Name == "Group B"

	if isGroupA {
		if rank == 1 {
			targetLabel = "SF1"
			targetCol = "team_a_id"
		} else {
			targetLabel = "SF2"
			targetCol = "team_b_id"
		}
	} else { // Group B
		if rank == 1 {
			targetLabel = "SF2"
			targetCol = "team_a_id"
		} else {
			targetLabel = "SF1"
			targetCol = "team_b_id"
		}
	}

	// 4. Update the Target Match
	for _, m := range koGroup.Matches {
		if m.Label == targetLabel {
			_, err := h.DB.NewUpdate().Model(m).Set(targetCol+" = ?", teamID).WherePK().Exec(ctx)
			return err
		}
	}

	return nil
}
