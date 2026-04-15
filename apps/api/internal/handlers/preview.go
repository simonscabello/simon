package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/services"
)

func PreviewRequest(svc *services.RequestPreviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		var uri requestIDURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}
		var q struct {
			EnvironmentID uint `form:"environment_id" binding:"required,min=1"`
		}
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "environment_id é obrigatório"})
			return
		}
		out, err := svc.Build(userID, uri.ID, q.EnvironmentID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "request ou ambiente não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível montar o preview"})
			return
		}
		c.JSON(http.StatusOK, out)
	}
}
