package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/services"
)

type keyValueParamBody struct {
	Key     string `json:"key" binding:"required,min=1,max=255"`
	Value   string `json:"value"`
	Enabled *bool  `json:"enabled"`
}

func CreateRequestHeader(svc *services.RequestHeaderService) gin.HandlerFunc {
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
		h, err := svc.Create(userID, uri.ID, body.Key, body.Value, body.Enabled)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "request não encontrado"})
				return
			}
			if errors.Is(err, services.ErrInvalidKey) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chave é obrigatória"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível criar o header"})
			return
		}
		c.JSON(http.StatusCreated, h)
	}
}

func ListRequestHeaders(svc *services.RequestHeaderService) gin.HandlerFunc {
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
		list, err := svc.List(userID, uri.ID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "request não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível listar headers"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func UpdateRequestHeader(svc *services.RequestHeaderService) gin.HandlerFunc {
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
		h, err := svc.Update(userID, uri.ID, body.Key, body.Value, body.Enabled)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "header não encontrado"})
				return
			}
			if errors.Is(err, services.ErrInvalidKey) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chave é obrigatória"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível atualizar o header"})
			return
		}
		c.JSON(http.StatusOK, h)
	}
}

func DeleteRequestHeader(svc *services.RequestHeaderService) gin.HandlerFunc {
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
				c.JSON(http.StatusNotFound, gin.H{"error": "header não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível remover o header"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
