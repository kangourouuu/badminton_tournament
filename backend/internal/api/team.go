package api

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"badminton_tournament/backend/internal/models"
)

type GenerateTeamsRequest struct {
	Pool string `json:"pool" binding:"required"`
}

func (h *Handler) GenerateTeams(c *gin.Context) {
	var req GenerateTeamsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. Fetch participants
	var participants []models.Participant
	err := h.DB.NewSelect().Model(&participants).
		Where("pool = ?", req.Pool).
		Scan(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(participants) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough participants to generate teams"})
		return
	}

	// 2. Shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})

	// 3. Pair them up
	var teams []*models.Team
	for i := 0; i < len(participants)-1; i += 2 {
		p1 := participants[i]
		p2 := participants[i+1]

		team := &models.Team{
			Player1ID: p1.ID,
			Player2ID: p2.ID,
			Pool:      req.Pool,
			Name:      p1.Name + " & " + p2.Name,
		}
		teams = append(teams, team)
	}

	// TODO: Handle odd one out (maybe return in response as 'waiting')
	
	// 4. Save to DB
	// Clear existing teams for this pool? Or just add? 
	// For simplicity, we just add. Admin can clear manually if needed or we assume fresh run.
	
	if len(teams) > 0 {
		_, err = h.DB.NewInsert().Model(&teams).Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"generated_count": len(teams),
		"teams":           teams,
		"waiting":         len(participants) % 2,
	})
}

func (h *Handler) ListTeams(c *gin.Context) {
	pool := c.Query("pool")
	var teams []models.Team
	query := h.DB.NewSelect().Model(&teams).Relation("Player1").Relation("Player2")

	if pool != "" {
		query.Where("tm.pool = ?", pool)
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

	if _, err := h.DB.NewInsert().Model(team).Exec(ctx); err != nil {
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
