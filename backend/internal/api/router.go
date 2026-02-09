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
	api.GET("/teams", h.ListTeams)
	api.GET("/groups", h.ListGroups)
	api.GET("/public/rules", h.GetRules)

	// Admin
	admin := api.Group("/")
	admin.Use(AuthMiddleware())
	{
		admin.POST("/teams/generate", h.GenerateTeams)
		admin.POST("/groups", h.CreateGroup)
		admin.POST("/matches/:id", h.UpdateMatch)
		admin.PUT("/admin/rules", h.UpdateRules)
	}
}
