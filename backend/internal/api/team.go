package api

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourname/badminton-manager/backend/internal/models"
)

type GenerateTeamsRequest struct {
	Pool string `json:"pool"`
}

func (h *Handler) HandleGenerateTeams(c *gin.Context) {
	var req GenerateTeamsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	pool := models.Pool(req.Pool)
	if pool != models.PoolMesoneer && pool != models.PoolLab {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pool"})
		return
	}

	// Fetch all participants for the pool
	var participants []models.Participant
	if err := h.DB.NewSelect().Model(&participants).Where("pool = ?", pool).Scan(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch participants"})
		return
	}

	if len(participants) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough participants to form a team"})
		return
	}

	// Shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})

	// Pair up
	var teams []models.Team
	for i := 0; i < len(participants); i += 2 {
		if i+1 < len(participants) {
			p1 := participants[i]
			p2 := participants[i+1]
			t := models.Team{
				Name:      p1.Name + " & " + p2.Name,
				Player1ID: p1.ID,
				Player2ID: p2.ID,
				Pool:      pool,
			}
			teams = append(teams, t)
		} else {
			// Handle odd number? For now, ignore or add to a "remainder" pool?
			// The prompt assumes pure pairing.
		}
	}

	if len(teams) > 0 {
		_, err := h.DB.NewInsert().Model(&teams).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save teams"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}
