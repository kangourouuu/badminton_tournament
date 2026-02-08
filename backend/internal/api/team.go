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
