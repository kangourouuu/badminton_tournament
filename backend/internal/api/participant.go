package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"badminton_tournament/backend/internal/models"
)

// Webhook for Google Form
// Expected format from Google Script
type GoogleFormRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Pool  string `json:"pool"`
}

func (h *Handler) HandleFormWebhook(c *gin.Context) {
	var req GoogleFormRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	participant := &models.Participant{
		Email: req.Email,
		Name:  req.Name,
		Pool:  req.Pool,
	}

	// Upsert: On conflict email, update name/pool
	_, err := h.DB.NewInsert().Model(participant).
		On("CONFLICT (email) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("pool = EXCLUDED.pool").
		Exec(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "id": participant.ID})
}

func (h *Handler) ListParticipants(c *gin.Context) {
	pool := c.Query("pool")

	var participants []models.Participant
	query := h.DB.NewSelect().Model(&participants)

	if pool != "" {
		query.Where("pool = ?", pool)
	}

	err := query.Order("name ASC").Scan(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, participants)
}
