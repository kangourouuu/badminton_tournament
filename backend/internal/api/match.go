package api

import (
	"context"
	"log"
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
	Status     string                 `json:"status"` // "finished" or empty
}

func (h *Handler) GetMatch(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	var match models.Match
	// Fetch match with related teams to ensure names are available
	err := h.DB.NewSelect().
		Model(&match).
		Relation("TeamA").
		Relation("TeamB").
		Where("id = ?", id).
		Scan(ctx)
		
	if err != nil {
		log.Printf("Error fetching match %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	c.JSON(http.StatusOK, match)
}

func (h *Handler) UpdateMatch(c *gin.Context) {
	id := c.Param("id")
	var req UpdateMatchRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation: If status is finished, winner must be set
	if req.Status == "finished" && req.WinnerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "WinnerID is required when status is finished"})
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
		// Propagate Winner
		// PRIORITIZE Group Stage Promotion (Winners/Decider) to enforce Cross-Over Logic
		// This must run BEFORE NextMatchWinID check to prevent legacy pointers from hijacking the route
		if match.Label == "Winners" { // M3 winner is Rank 1
			log.Printf("[Auto-Propagation] Promoting Group Rank 1 (Winner %s) to Knockout", req.WinnerID)
			h.promoteToKnockout(ctx, match.GroupID, 1, req.WinnerID)
		} else if match.Label == "Decider" { // M5 winner is Rank 2
			log.Printf("[Auto-Propagation] Promoting Group Rank 2 (Decider Winner %s) to Knockout", req.WinnerID)
			if err := h.promoteToKnockout(ctx, match.GroupID, 2, req.WinnerID); err != nil {
				log.Printf("[Auto-Propagation] ERROR promoting decider: %v", err)
			}
		} else if match.NextMatchWinID != uuid.Nil {
			log.Printf("[Auto-Propagation] Propagating WINNER %s to Match %s (Source: %s)", req.WinnerID, match.NextMatchWinID, match.Label)
			h.propagateToMatch(ctx, match.NextMatchWinID, req.WinnerID, match.Label, "win")
		} else {
             // Fallback for others
		}

		// Propagate Loser
		if match.NextMatchLoseID != uuid.Nil && loserID != uuid.Nil {
			log.Printf("[Auto-Propagation] Propagating LOSER %s to Match %s (Source: %s)", loserID, match.NextMatchLoseID, match.Label)
			// Ensure "lose" outcome is passed for Bronze logic
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
		log.Printf("PROMOTION SUCCESS: Pushed Player %v to Match ID %v", teamID, target.ID)
		_, err := h.DB.NewUpdate().Model(&target).Set(col+" = ?", teamID).WherePK().Exec(ctx)
		return err
	}
	return nil
}

func (h *Handler) promoteToKnockout(ctx context.Context, groupID uuid.UUID, rank int, teamID uuid.UUID) error {
	log.Printf("DEBUG: promoteToKnockout called for Group %v, Rank %d, Team %v", groupID, rank, teamID)

	var group models.Group
	if err := h.DB.NewSelect().Model(&group).Where("id = ?", groupID).Scan(ctx); err != nil {
		log.Printf("PROMOTION ERROR: Source Group %v not found: %v", groupID, err)
		return err
	}
	log.Printf("DEBUG: Promotion Source Group: Name='%s', Pool='%s'", group.Name, group.Pool)

	// 2. Find Knockout Group (or Auto-Generate)
	var koGroup models.Group
	if err := h.DB.NewSelect().Model(&koGroup).Where("name = ?", "KNOCKOUT").Relation("Matches").Scan(ctx); err != nil {
		// Try looser search
		if err2 := h.DB.NewSelect().Model(&koGroup).Where("name ILIKE ?", "knockout%").Relation("Matches").Scan(ctx); err2 != nil {
			log.Printf("PROMOTION NOTICE: Knockout Stage not found. Attempting Auto-Generation...")
			
			// Auto-Generate
			newKoGroup, errGen := h.EnsureKnockoutStage(ctx, group.TournamentID)
			if errGen != nil {
				log.Printf("PROMOTION ERROR: Failed to auto-generate Knockout Stage: %v", errGen)
				return errGen 
			}
			koGroup = *newKoGroup
			log.Printf("PROMOTION SUCCESS: Auto-Generated Knockout Stage (ID: %s)", koGroup.ID)
		}
	}
	log.Printf("DEBUG: Found Knockout Group %v with %d matches", koGroup.ID, len(koGroup.Matches))
	for _, m := range koGroup.Matches {
		log.Printf("DEBUG: Available Match in KO Group: ID=%s Label=%s", m.ID, m.Label)
	}
	log.Printf("DEBUG: Found Knockout Group %v with %d matches", koGroup.ID, len(koGroup.Matches))

	var targetLabel string
	var targetCol string

	// MASTERPLAN Macro-Flow Rule: Cross-Over Semi-Finals
	// Rank 1 Group A (Mesoneer) vs Rank 2 Group B (Lab) -> SF1
	// Rank 1 Group B (Lab)      vs Rank 2 Group A (Mesoneer) -> SF2

	isPoolA := group.Pool == "Mesoneer"
	isPoolA := group.Pool == "Mesoneer"
	if isPoolA {
		if rank == 1 {
			// TICKET 1: Group M Winner -> SF1 (Slot 1)
			targetLabel = "SF1"
			targetCol = "team_a_id"
		} else { 
			// TICKET 2: Group M Runner-up -> SF2 (Slot 2) [Cross-over]
			targetLabel = "SF2"
			targetCol = "team_b_id"
		}
	} else { // Pool B (Lab)
		if rank == 1 {
			// TICKET 1: Group L Winner -> SF2 (Slot 1)
			targetLabel = "SF2"
			targetCol = "team_a_id"
		} else { 
			// TICKET 2: Group L Runner-up -> SF1 (Slot 2) [Cross-over]
			targetLabel = "SF1"
			targetCol = "team_b_id"
		}
	}

	// 1. Search in Loaded Relation
	for _, m := range koGroup.Matches {
		if m.Label == targetLabel {
			log.Printf("DEBUG: Found Target Match %s (ID: %s) in relation", targetLabel, m.ID)
			res, err := h.DB.NewUpdate().Model(m).Set(targetCol+" = ?", teamID).WherePK().Exec(ctx)
			if err != nil {
				log.Printf("PROMOTION ERROR: DB Update failed: %v", err)
				return err
			}
			rows, _ := res.RowsAffected()
			log.Printf("PROMOTION SUCCESS: Updated %s with Team %s. Rows Affected: %d", targetLabel, teamID, rows)
			return nil
		}
	}

	// 2. Fallback: Direct Query (if koGroup.Matches relation was empty/incomplete)
	log.Printf("PROMOTION WARNING: Target %s not found in relation, trying direct query", targetLabel)
	var targetMatch models.Match
	// Note: Explicitly selecting ID to ensure we have a valid PK for update
	if err := h.DB.NewSelect().Model(&targetMatch).Where("group_id = ? AND label = ?", koGroup.ID, targetLabel).Scan(ctx); err == nil {
		log.Printf("DEBUG: Found Target Match %s (ID: %s) via direct query", targetLabel, targetMatch.ID)
		res, err := h.DB.NewUpdate().Model(&targetMatch).Set(targetCol+" = ?", teamID).WherePK().Exec(ctx)
		if err != nil {
			log.Printf("PROMOTION ERROR: DB Update failed (fallback): %v", err)
			return err
		}
		rows, _ := res.RowsAffected()
		log.Printf("PROMOTION SUCCESS (FALLBACK): Updated %s with Team %s. Rows Affected: %d", targetLabel, teamID, rows)
		return nil
	}
	
	log.Printf("PROMOTION ERROR: Target Match %s not found in DB for progression", targetLabel)
	return nil
}
