package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"badminton_tournament/backend/internal/models"
)

type UpdateRulesRequest struct {
	Content string `json:"content" binding:"required"`
}

// GetRules fetches the latest rule content
// GET /api/public/rules
func (h *Handler) GetRules(c *gin.Context) {
	var rule models.Rule
	// Just get the last updated one or any one, since we likely have only one singleton rule
	err := h.DB.NewSelect().Model(&rule).Order("updated_at DESC").Limit(1).Scan(c.Request.Context())
	
	if err != nil {
		// If no rows, return empty default rule
		c.JSON(http.StatusOK, models.Rule{Content: "No rules defined yet."})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// UpdateRules updates or creates the rule content
// PUT /api/admin/rules
func (h *Handler) UpdateRules(c *gin.Context) {
	var req UpdateRulesRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. Try to find existing rule
	var rule models.Rule
	exists, _ := h.DB.NewSelect().Model(&rule).Exists(ctx)

	if exists {
		// Update existing (fetch ID first or just update all? Let's fetch last one)
		h.DB.NewSelect().Model(&rule).Order("updated_at DESC").Limit(1).Scan(ctx)
		rule.Content = req.Content
		rule.UpdatedAt = time.Now()
		_, err := h.DB.NewUpdate().Model(&rule).WherePK().Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Create new
		rule = models.Rule{
			Content: req.Content,
		}
		_, err := h.DB.NewInsert().Model(&rule).Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, rule)
}
