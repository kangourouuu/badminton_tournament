package middleware

import (
	"log"
	"net/http"
	"strings"

	"badminton_tournament/backend/internal/auth"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware() gin.HandlerFunc {
	// Initialize Casbin enforcer
	e, err := casbin.NewEnforcer("auth/model.conf", "auth/policy.csv")
	if err != nil {
		log.Fatalf("Casbin init failure: %v", err)
	}

	return func(c *gin.Context) {
		// Get Token from Header
		authHeader := c.GetHeader("Authorization")
		role := "viewer" // Default role

		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := auth.ValidateToken(tokenString)
			if err == nil {
				role = claims.Role
			}
		}

		// Enforce Policy
		obj := c.Request.URL.Path
		act := c.Request.Method

		allowed, err := e.Enforce(role, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Authorization error"})
			return
		}

		if allowed {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
	}
}
