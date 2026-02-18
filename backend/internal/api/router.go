package api

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type Handler struct {
	DB *bun.DB
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	
	// Auth
	api.POST("/auth/login", h.Login)
	api.POST("/webhooks/form", h.HandleFormWebhook)

	// Public
	api.GET("/participants", h.ListParticipants)
	api.POST("/participants", h.HandleFormWebhook) // Endpoint for Google Form Script
	api.GET("/teams", h.ListTeams)
	api.GET("/groups", h.ListGroups)
	api.GET("/matches/:id", h.GetMatch)
	api.GET("/public/rules", h.GetRules)

	// Admin
	admin := api.Group("/")
	admin.Use(AuthMiddleware())
	{

		admin.POST("/teams", h.CreateTeam)
		admin.POST("/teams/auto-pair", h.AutoPairTeams)
		admin.PUT("/teams/:id", h.UpdateTeam)
		admin.DELETE("/teams/:id", h.DeleteTeam)
		admin.POST("/groups", h.CreateGroup)
		admin.POST("/groups/auto-generate", h.AutoGenerateGroups)
		admin.POST("/matches/:id", h.UpdateMatch)
		admin.POST("/tournaments/knockout", h.GenerateKnockout)
		admin.PUT("/admin/rules", h.UpdateRules)
	}
}
