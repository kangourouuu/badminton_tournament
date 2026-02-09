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

	if col != "" {
		_, err = h.DB.NewUpdate().Model(&target).Column(col).WherePK().Exec(ctx)
		return err
	}
	return nil
}
