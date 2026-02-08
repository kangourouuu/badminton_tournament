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
			h.propagateToMatch(ctx, match.NextMatchWinID, req.WinnerID)
		}

		// Propagate Loser
		if match.NextMatchLoseID != uuid.Nil && loserID != uuid.Nil {
			h.propagateToMatch(ctx, match.NextMatchLoseID, loserID)
		}
	}

	c.JSON(http.StatusOK, match)
}

func (h *Handler) propagateToMatch(ctx context.Context, matchID, teamID uuid.UUID) error {
	// Find target match
	var target models.Match
	err := h.DB.NewSelect().Model(&target).Where("id = ?", matchID).Scan(ctx)
	if err != nil {
		return err
	}

	// Check which slot is empty or if we are overwriting
	// Logic: If TeamA is empty, fill it. Else if TeamB is empty, fill it.
	// If both full, maybe we are correcting a mistake? For now, fill first empty.
	
	cols := []string{}

	if target.TeamAID == uuid.Nil {
		target.TeamAID = teamID
		cols = append(cols, "team_a_id")
	} else if target.TeamBID == uuid.Nil {
		target.TeamBID = teamID
		cols = append(cols, "team_b_id")
	} else {
		// Both full. In a real app, maybe check if one of them was the previous winner/loser from this source?
		// For MVP, we stick to "First Empty Slot".
		// Or maybe force overwrite? Let's assume the flow is clean.
		// If Admin made a mistake and re-updates, we might have an issue.
		// Detailed logic:
		// We need to know WHICH slot came from THIS match. 
		// But 'Match' struct doesn't strictly store "SourceMatchID".
		// Simple fix: If target has TeamA, check if TeamA came from here? Hard to track.
		// MVP: Just fill TeamA if empty, else TeamB.
		return nil 
	}

	if len(cols) > 0 {
		_, err = h.DB.NewUpdate().Model(&target).Column(cols...).WherePK().Exec(ctx)
		return err
	}
	return nil
}
