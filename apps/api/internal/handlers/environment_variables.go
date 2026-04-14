package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/services"
)

func CreateEnvironmentVariable(svc *services.EnvironmentVariableService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		var uri envIDURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}
		var body keyValueParamBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		v, err := svc.Create(userID, uri.ID, body.Key, body.Value, body.Enabled)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "ambiente não encontrado"})
				return
			}
			if errors.Is(err, services.ErrInvalidKey) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chave é obrigatória"})
				return
			}
			if errors.Is(err, services.ErrDuplicateKey) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chave já existe neste ambiente"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível criar a variável"})
			return
		}
		c.JSON(http.StatusCreated, v)
	}
}

func ListEnvironmentVariables(svc *services.EnvironmentVariableService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		var uri envIDURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}
		list, err := svc.List(userID, uri.ID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "ambiente não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível listar variáveis"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func UpdateEnvironmentVariable(svc *services.EnvironmentVariableService) gin.HandlerFunc {
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
		var body keyValueParamBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		v, err := svc.Update(userID, uri.ID, body.Key, body.Value, body.Enabled)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "variável não encontrada"})
				return
			}
			if errors.Is(err, services.ErrInvalidKey) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chave é obrigatória"})
				return
			}
			if errors.Is(err, services.ErrDuplicateKey) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chave já existe neste ambiente"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível atualizar a variável"})
			return
		}
		c.JSON(http.StatusOK, v)
	}
}

func DeleteEnvironmentVariable(svc *services.EnvironmentVariableService) gin.HandlerFunc {
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
		if err := svc.Delete(userID, uri.ID); err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "variável não encontrada"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível remover a variável"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
