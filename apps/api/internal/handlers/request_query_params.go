package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/services"
)

func CreateRequestQueryParam(svc *services.RequestQueryParamService) gin.HandlerFunc {
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
		q, err := svc.Create(userID, uri.ID, body.Key, body.Value, body.Enabled)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "request não encontrado"})
				return
			}
			if errors.Is(err, services.ErrInvalidKey) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chave é obrigatória"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível criar o query param"})
			return
		}
		c.JSON(http.StatusCreated, q)
	}
}

func ListRequestQueryParams(svc *services.RequestQueryParamService) gin.HandlerFunc {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível listar query params"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func UpdateRequestQueryParam(svc *services.RequestQueryParamService) gin.HandlerFunc {
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
		q, err := svc.Update(userID, uri.ID, body.Key, body.Value, body.Enabled)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "query param não encontrado"})
				return
			}
			if errors.Is(err, services.ErrInvalidKey) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chave é obrigatória"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível atualizar o query param"})
			return
		}
		c.JSON(http.StatusOK, q)
	}
}

func DeleteRequestQueryParam(svc *services.RequestQueryParamService) gin.HandlerFunc {
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
				c.JSON(http.StatusNotFound, gin.H{"error": "query param não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível remover o query param"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
