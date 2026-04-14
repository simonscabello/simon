package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"simon/apps/api/internal/middlewares"
	"simon/apps/api/internal/repositories"
)

func Me(db *gorm.DB) gin.HandlerFunc {
	repo := repositories.NewUserRepository(db)
	return func(c *gin.Context) {
		v, ok := c.Get(middlewares.ContextUserIDKey)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		userID, ok := v.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		user, err := repo.FindByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível carregar o usuário"})
			return
		}
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
			return
		}
		user.Password = ""
		c.JSON(http.StatusOK, user)
	}
}
