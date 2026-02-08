package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/badminton-manager/backend/internal/models"
)

type WebhookRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Pool  string `json:"pool"` // "Mesoneer" or "Lab"
}

func (h *Handler) HandleWebhookForm(c *gin.Context) {
	var req WebhookRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Map generic string to Enum
	var pool models.Pool
	switch req.Pool {
	case "Mesoneer":
		pool = models.PoolMesoneer
	case "Lab":
		pool = models.PoolLab
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pool"})
		return
	}

	p := &models.Participant{
		Email: req.Email,
		Name:  req.Name,
		Pool:  pool,
	}

	// Upsert: On Conflict (Email) Do Update
	_, err := h.DB.NewInsert().Model(p).
		On("CONFLICT (email) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("pool = EXCLUDED.pool").
		Exec(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save participant: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "participant": p})
}
