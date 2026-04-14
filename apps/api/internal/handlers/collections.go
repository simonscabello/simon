package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/middlewares"
	"simon/apps/api/internal/services"
)

type collectionCreateBody struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

type collectionUpdateBody struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

type collectionIDURI struct {
	ID uint `uri:"id" binding:"required"`
}

func CreateCollection(svc *services.CollectionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		var body collectionCreateBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		col, err := svc.Create(userID, body.Name)
		if err != nil {
			if errors.Is(err, services.ErrInvalidName) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "nome é obrigatório"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível criar a collection"})
			return
		}
		c.JSON(http.StatusCreated, col)
	}
}

func ListCollections(svc *services.CollectionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		list, err := svc.List(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível listar collections"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func UpdateCollection(svc *services.CollectionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		var uri collectionIDURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}
		var body collectionUpdateBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		col, err := svc.Update(userID, uri.ID, body.Name)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "collection não encontrada"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível atualizar a collection"})
			return
		}
		c.JSON(http.StatusOK, col)
	}
}

func DeleteCollection(svc *services.CollectionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		var uri collectionIDURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}
		if err := svc.Delete(userID, uri.ID); err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "collection não encontrada"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível remover a collection"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func userIDFromContext(c *gin.Context) (uint, bool) {
	v, ok := c.Get(middlewares.ContextUserIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
