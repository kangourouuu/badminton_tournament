package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"badminton_tournament/backend/internal/models"
)



func (h *Handler) ListTeams(c *gin.Context) {
	pool := c.Query("pool")
	available := c.Query("available") == "true"

	var teams []models.Team
	query := h.DB.NewSelect().Model(&teams).Relation("Player1").Relation("Player2")

	if pool != "" {
		query.Where("tm.pool = ?", pool)
	}

	if available {
		// Filter out teams that are already in any match
		// Subquery: SELECT team_a_id FROM matches UNION SELECT team_b_id FROM matches
		// Bun doesn't support UNION easily in subquery builder sometimes, so we can do WHERE NOT EXISTS or NOT IN
		// Simpler: WHERE id NOT IN (SELECT team_a_id FROM matches) AND id NOT IN (SELECT team_b_id FROM matches)
		
		query.Where("id NOT IN (SELECT team_a_id FROM matches WHERE team_a_id IS NOT NULL)")
		query.Where("id NOT IN (SELECT team_b_id FROM matches WHERE team_b_id IS NOT NULL)")
	}

	err := query.Order("tm.created_at DESC").Scan(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// CreateTeamRequest for manual team creation
type CreateTeamRequest struct {
	Player1ID string `json:"player1_id" binding:"required"`
	Player2ID string `json:"player2_id" binding:"required"`
}

// CreateTeam - Manual creation with strict pool validation
func (h *Handler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. Fetch Players
	var p1, p2 models.Participant
	if err := h.DB.NewSelect().Model(&p1).Where("id = ?", req.Player1ID).Scan(ctx); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player 1 not found"})
		return
	}
	if err := h.DB.NewSelect().Model(&p2).Where("id = ?", req.Player2ID).Scan(ctx); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player 2 not found"})
		return
	}

	// 2. Validate Pool
	if p1.Pool != p2.Pool {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Players must be from the same pool"})
		return
	}

	// 3. Validate Availability: Check if players are already in a team
	count, _ := h.DB.NewSelect().Model((*models.Team)(nil)).
		Where("player1_id IN (?) OR player2_id IN (?)", bun.In([]string{req.Player1ID, req.Player2ID}), bun.In([]string{req.Player1ID, req.Player2ID})).
		Count(ctx)
	
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or both players are already in a team"})
		return
	}

	// 4. Create Team
	team := &models.Team{
		Player1ID: p1.ID,
		Player2ID: p2.ID,
		Pool:      p1.Pool, // Inherit pool
		Name:      p1.Name + " & " + p2.Name,
	}

	if _, err := h.DB.NewInsert().Model(team).Returning("*").Exec(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// UpdateTeamRequest - Swap players
type UpdateTeamRequest struct {
	Player1ID string `json:"player1_id"`
	Player2ID string `json:"player2_id"`
}

// UpdateTeam - Edit composition
func (h *Handler) UpdateTeam(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTeamRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	var team models.Team
	if err := h.DB.NewSelect().Model(&team).Where("id = ?", id).Scan(ctx); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// If P1 passed, update
	if req.Player1ID != "" {
		var p models.Participant
		if err := h.DB.NewSelect().Model(&p).Where("id = ?", req.Player1ID).Scan(ctx); err == nil {
			if p.Pool != team.Pool {
				c.JSON(http.StatusBadRequest, gin.H{"error": "New Player 1 is from wrong pool"})
				return
			}
			team.Player1ID = p.ID
		}
	}
	
	// If P2 passed, update
	if req.Player2ID != "" {
		var p models.Participant
		if err := h.DB.NewSelect().Model(&p).Where("id = ?", req.Player2ID).Scan(ctx); err == nil {
			if p.Pool != team.Pool {
				c.JSON(http.StatusBadRequest, gin.H{"error": "New Player 2 is from wrong pool"})
				return
			}
			team.Player2ID = p.ID
		}
	}

	// Update Name
	// Ideally we re-fetch both to concat names, but let's do a quick lazy fetch or just leave it.
	// Correct way: Fetch both current IDs to rebuild name.
	var p1, p2 models.Participant
	h.DB.NewSelect().Model(&p1).Where("id = ?", team.Player1ID).Scan(ctx)
	h.DB.NewSelect().Model(&p2).Where("id = ?", team.Player2ID).Scan(ctx)
	team.Name = p1.Name + " & " + p2.Name

	if _, err := h.DB.NewUpdate().Model(&team).WherePK().Exec(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// DeleteTeam - Disband
func (h *Handler) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	// Check for finished matches?
	count, _ := h.DB.NewSelect().Model((*models.Match)(nil)).
		Where("team_a_id = ? OR team_b_id = ?", id, id).
		Where("winner_id IS NOT NULL").
		Count(ctx)
	
	if count > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot disband team that has finished matches"})
		return
	}

	if _, err := h.DB.NewDelete().Model((*models.Team)(nil)).Where("id = ?", id).Exec(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team disbanded"})
}
