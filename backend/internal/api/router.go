package api

import (
	"badminton_tournament/backend/internal/auth"
	"badminton_tournament/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Public
	r.POST("/api/auth/login", auth.LoginHandler)
	r.GET("/api/bracket", h.HandleGetBracket)
	r.GET("/api/teams", h.HandleGetTeams) // Admin uses this too, could be protected but GET is safe-ish
	r.POST("/api/webhooks/form", h.HandleWebhookForm)

	// Protected
	api := r.Group("/api")
	api.Use(middleware.CasbinMiddleware())
	{
		api.POST("/teams/generate", h.HandleGenerateTeams)
		api.POST("/matches/:id/result", h.HandleMatchResult)
		api.POST("/bracket/generate", h.HandleGenerateBracket)
	}
}
