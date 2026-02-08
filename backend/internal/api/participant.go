package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/badminton-manager/backend/internal/models"
)

type WebhookRequest struct {
	Name string `json:"name"`
	Pool string `json:"pool"` // "Mesoneer" or "Lab"
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
		Name: req.Name,
		Pool: pool,
	}

	_, err := h.DB.NewInsert().Model(p).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save participant: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "participant": p})
}
