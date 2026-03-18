package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"badminton_tournament/backend/internal/models"
)

// Webhook for Google Form
// Expected format from Google Script
// Webhook for Google Form
// Updated format: { "name": "...", "group": "...", "categories": ["..."], "available_dates": ["..."] }
type GoogleFormRequest struct {
	Name           string   `json:"name"`
	Group          string   `json:"group"` // Maps to Pool
	Categories     []string `json:"categories"`
	AvailableDates []string `json:"available_dates"`
	Gender         string   `json:"gender"`
	Source         string   `json:"source"`
	Status         string   `json:"status"`
}

func (h *Handler) HandleFormWebhook(c *gin.Context) {
	var req GoogleFormRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Auto-generate pseudo-email removed. Usage Name as unique key.
	
	participant := &models.Participant{
		Name:           req.Name,
		Pool:           req.Group, // Google Form "group" -> DB "pool"
		Categories:     req.Categories,
		AvailableDates: req.AvailableDates,
		Gender:         req.Gender,
		Source:         req.Source,
		Status:         req.Status,
	}

	// Upsert: On conflict name, update pool/categories/available_dates
	_, err := h.DB.NewInsert().Model(participant).
		On("CONFLICT (name) DO UPDATE").
		Set("pool = EXCLUDED.pool").
		Set("categories = EXCLUDED.categories").
		Set("available_dates = EXCLUDED.available_dates").
		Set("gender = EXCLUDED.gender").
		Set("source = EXCLUDED.source").
		Set("status = EXCLUDED.status").
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
