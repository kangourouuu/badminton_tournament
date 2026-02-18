package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"badminton_tournament/backend/internal/models"
)

// Webhook for Google Form
// Expected format from Google Script
// Webhook for Google Form
// Updated format: { "name": "...", "group": "...", "partner_request": "..." }
type GoogleFormRequest struct {
	Name           string `json:"name"`
	Group          string `json:"group"`           // Maps to Pool
	PartnerRequest string `json:"partner_request"` // Desired Teammate
}

func (h *Handler) HandleFormWebhook(c *gin.Context) {
	var req GoogleFormRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Auto-generate pseudo-email since it's removed from form but required by DB
	// Format: name_sanitized@placeholder.badminton
	// This relies on Name being unique enough, or we accept overwrite if Name is identical.
	sanitized := models.SanitizeNameForEmail(req.Name)
	pseudoEmail := sanitized + "@placeholder.badminton"

	participant := &models.Participant{
		Email:          pseudoEmail,
		Name:           req.Name,
		Pool:           req.Group, // Google Form "group" -> DB "pool"
		PartnerRequest: req.PartnerRequest,
	}

	// Upsert: On conflict email (which is driven by name now), update info
	_, err := h.DB.NewInsert().Model(participant).
		On("CONFLICT (email) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("pool = EXCLUDED.pool").
		Set("partner_request = EXCLUDED.partner_request").
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
