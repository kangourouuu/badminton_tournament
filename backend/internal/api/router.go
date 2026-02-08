package api

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/webhooks/form", h.HandleWebhookForm)
		api.POST("/teams/generate", h.HandleGenerateTeams)
		api.POST("/matches/:id/result", h.HandleMatchResult)
		api.POST("/bracket/generate", h.HandleGenerateBracket)
		api.GET("/bracket", h.HandleGetBracket)
	}
}
