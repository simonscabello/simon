package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/auth"
)

const ContextUserIDKey = "user_id"

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token ausente ou inválido"})
			c.Abort()
			return
		}
		raw := strings.TrimSpace(h[7:])
		if raw == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token ausente ou inválido"})
			c.Abort()
			return
		}
		claims, err := auth.ParseAccessToken(raw, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token ausente ou inválido"})
			c.Abort()
			return
		}
		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}
