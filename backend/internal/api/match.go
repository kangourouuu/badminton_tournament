package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourname/badminton-manager/backend/internal/models"
)

type GenerateBracketRequest struct {
	TeamIDs []int64 `json:"team_ids"`
	Pool    string  `json:"pool"`
}

func (h *Handler) HandleGenerateBracket(c *gin.Context) {
	var req GenerateBracketRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.TeamIDs) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exactly 4 teams required for GSL group"})
		return
	}

	ctx := c.Request.Context()

	// Create 5 matches
	// M1: T1 vs T2
	// M2: T3 vs T4
	// M3: Winner M1 vs Winner M2 (Winners Match)
	// M4: Loser M1 vs Loser M2 (Losers Match)
	// M5: Loser M3 vs Winner M4 (Decider Match)

	// In reverse order to get IDs for linking? Or update later.
	// Let's create empty matches first or use a transaction.
	
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback()

	// Helper to create match
	createMatch := func(desc string, tA, tB int64) *models.Match {
		m := &models.Match{
			Description: desc,
			Pool:        models.Pool(req.Pool),
			TeamAID:     tA,
			TeamBID:     tB,
		}
		if _, err := tx.NewInsert().Model(m).Exec(ctx); err != nil {
			panic(err) // Handle better in real code
		}
		return m
	}

	// We need to create M5 first to link M3/M4? No, we can update.
	// Let's create all 5 matches first.
	m1 := createMatch("Opening Match 1", req.TeamIDs[0], req.TeamIDs[1])
	m2 := createMatch("Opening Match 2", req.TeamIDs[2], req.TeamIDs[3])
	m3 := createMatch("Winners Match", 0, 0)
	m4 := createMatch("Elimination Match", 0, 0)
	m5 := createMatch("Decider Match", 0, 0)

	// Link M1 -> Winner to M3 (A), Loser to M4 (A)
	m1.NextMatchID = m3.ID
	m1.NextMatchSlot = "A"
	m1.LoserNextMatchID = m4.ID
	m1.LoserNextMatchSlot = "A"
	tx.NewUpdate().Model(m1).WherePK().Exec(ctx)

	// Link M2 -> Winner to M3 (B), Loser to M4 (B)
	m2.NextMatchID = m3.ID
	m2.NextMatchSlot = "B"
	m2.LoserNextMatchID = m4.ID
	m2.LoserNextMatchSlot = "B"
	tx.NewUpdate().Model(m2).WherePK().Exec(ctx)

	// Link M3 -> Winner (Group Winner), Loser to M5 (A)
	m3.LoserNextMatchID = m5.ID
	m3.LoserNextMatchSlot = "A"
	tx.NewUpdate().Model(m3).WherePK().Exec(ctx)

	// Link M4 -> Winner to M5 (B), Loser (Eliminated)
	m4.NextMatchID = m5.ID
	m4.NextMatchSlot = "B"
	tx.NewUpdate().Model(m4).WherePK().Exec(ctx)
	
	// M5 -> Winner (Group Runner-up)

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bracket generated", "matches": []int64{m1.ID, m2.ID, m3.ID, m4.ID, m5.ID}})
}

type MatchResultRequest struct {
	ScoreA   int    `json:"score_a"`
	ScoreB   int    `json:"score_b"`
	VideoURL string `json:"video_url"`
}

func (h *Handler) HandleMatchResult(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req MatchResultRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := c.Request.Context()
	
	var match models.Match
	if err := h.DB.NewSelect().Model(&match).Where("id = ?", id).Scan(ctx); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Update score
	match.ScoreA = req.ScoreA
	match.ScoreB = req.ScoreB
	match.VideoURL = req.VideoURL

	// Determine winner
	var winnerID, loserID int64
	if req.ScoreA > req.ScoreB {
		winnerID = match.TeamAID
		loserID = match.TeamBID
	} else {
		winnerID = match.TeamBID
		loserID = match.TeamAID
	}
	match.WinnerID = winnerID

	if _, err := h.DB.NewUpdate().Model(&match).WherePK().Exec(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update match"})
		return
	}

	// Advance Winner
	if match.NextMatchID != 0 && winnerID != 0 {
		var nextMatch models.Match
		if err := h.DB.NewSelect().Model(&nextMatch).Where("id = ?", match.NextMatchID).Scan(ctx); err == nil {
			if match.NextMatchSlot == "A" {
				nextMatch.TeamAID = winnerID
			} else {
				nextMatch.TeamBID = winnerID
			}
			h.DB.NewUpdate().Model(&nextMatch).WherePK().Exec(ctx)
		}
	}

	// Advance Loser (Bracket specific)
	if match.LoserNextMatchID != 0 && loserID != 0 {
		var nextMatch models.Match
		if err := h.DB.NewSelect().Model(&nextMatch).Where("id = ?", match.LoserNextMatchID).Scan(ctx); err == nil {
			if match.LoserNextMatchSlot == "A" {
				nextMatch.TeamAID = loserID
			} else {
				nextMatch.TeamBID = loserID
			}
			h.DB.NewUpdate().Model(&nextMatch).WherePK().Exec(ctx)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "winner_id": winnerID})
}

func (h *Handler) HandleGetBracket(c *gin.Context) {
	pool := c.Query("pool")
	var matches []models.Match
	
	query := h.DB.NewSelect().Model(&matches).Relation("TeamA").Relation("TeamB")
	if pool != "" {
		query.Where("m.pool = ?", pool)
	}
	
	if err := query.Order("id ASC").Scan(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch matches"})
		return
	}

	c.JSON(http.StatusOK, matches)
}
