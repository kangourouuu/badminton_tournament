package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"badminton_tournament/backend/internal/models"
)

type UpdateMatchRequest struct {
	WinnerID   uuid.UUID              `json:"winner_id"`
	Score      string                 `json:"score"`
	SetsDetail map[string]interface{} `json:"sets_detail"`
	VideoURL   string                 `json:"video_url"`
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
	match.SetsDetail = req.SetsDetail
	match.VideoURL = req.VideoURL
	
	_, err = h.DB.NewUpdate().Model(&match).Column("winner_id", "score", "sets_detail", "video_url").WherePK().Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Auto-Propagation
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
		} else {
			// Check for qualification if no NextMatchWinID (M3 or M5)
			if match.Label == "Winners" { // M3 winner is Rank 1
				h.promoteToKnockout(ctx, match.GroupID, 1, req.WinnerID)
			} else if match.Label == "Decider" { // M5 winner is Rank 2
				h.promoteToKnockout(ctx, match.GroupID, 2, req.WinnerID)
			} else if match.Label == "Final" || match.Label == "Bronze" {
				// Champion decided, potentially update tournament status or just stay finished
			}
		}

		// Propagate Loser
		if match.NextMatchLoseID != uuid.Nil && loserID != uuid.Nil {
			h.propagateToMatch(ctx, match.NextMatchLoseID, loserID, match.Label, "lose")
		}
	}

	c.JSON(http.StatusOK, match)
}

func (h *Handler) propagateToMatch(ctx context.Context, targetID, teamID uuid.UUID, sourceLabel, outcome string) error {
	var target models.Match
	if err := h.DB.NewSelect().Model(&target).Where("id = ?", targetID).Scan(ctx); err != nil {
		return err
	}

	col := ""
	// MASTERPLAN Logic for GSL
	switch target.Label {
	case "Winners": // M3
		if sourceLabel == "M1" {
			col = "team_a_id"
		} else if sourceLabel == "M2" {
			col = "team_b_id"
		}
	case "Losers": // M4
		if sourceLabel == "M1" {
			col = "team_a_id"
		} else if sourceLabel == "M2" {
			col = "team_b_id"
		}
	case "Decider": // M5
		if sourceLabel == "Winners" { // Loser of M3 goes to Slot 1
			col = "team_a_id"
		} else if sourceLabel == "Losers" { // Winner of M4 goes to Slot 2
			col = "team_b_id"
		}
	case "Final": // Knockout Final
		// Source could be SF1 or SF2
		if sourceLabel == "SF1" {
			col = "team_a_id"
		} else if sourceLabel == "SF2" {
			col = "team_b_id"
		}
	case "Bronze":
		if sourceLabel == "SF1" {
			col = "team_a_id"
		} else if sourceLabel == "SF2" {
			col = "team_b_id"
		}
	}

	// Fallback to "First Empty" if label specific logic not matched
	if col == "" {
		if target.TeamAID == uuid.Nil {
			col = "team_a_id"
		} else if target.TeamBID == uuid.Nil {
			col = "team_b_id"
		}
	}

	if col != "" {
		_, err := h.DB.NewUpdate().Model(&target).Set(col+" = ?", teamID).WherePK().Exec(ctx)
		return err
	}
	return nil
}

func (h *Handler) promoteToKnockout(ctx context.Context, groupID uuid.UUID, rank int, teamID uuid.UUID) error {
	var group models.Group
	if err := h.DB.NewSelect().Model(&group).Where("id = ?", groupID).Scan(ctx); err != nil {
		return err
	}

	var koGroup models.Group
	if err := h.DB.NewSelect().Model(&koGroup).Where("name = ?", "KNOCKOUT").Relation("Matches").Scan(ctx); err != nil {
		return nil // Knockout bracket not yet created
	}

	var targetLabel string
	var targetCol string

	// MASTERPLAN Macro-Flow Rule: Cross-Over Semi-Finals
	// Rank 1 Group A (Mesoneer) vs Rank 2 Group B (Lab) -> SF1
	// Rank 1 Group B (Lab)      vs Rank 2 Group A (Mesoneer) -> SF2

	isPoolA := group.Pool == "Mesoneer"
	if isPoolA {
		if rank == 1 {
			targetLabel = "SF1"
			targetCol = "team_a_id"
		} else { // Runner-up Group A goes to SF2
			targetLabel = "SF2"
			targetCol = "team_b_id"
		}
	} else { // Pool B (Lab)
		if rank == 1 {
			targetLabel = "SF2"
			targetCol = "team_a_id"
		} else { // Runner-up Group B goes to SF1
			targetLabel = "SF1"
			targetCol = "team_b_id"
		}
	}

	for _, m := range koGroup.Matches {
		if m.Label == targetLabel {
			_, err := h.DB.NewUpdate().Model(m).Set(targetCol+" = ?", teamID).WherePK().Exec(ctx)
			return err
		}
	}
	return nil
}
