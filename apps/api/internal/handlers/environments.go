package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/services"
)

type environmentNameBody struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

type envIDURI struct {
	ID uint `uri:"id" binding:"required"`
}

func CreateEnvironment(svc *services.EnvironmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		var body environmentNameBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		e, err := svc.Create(userID, body.Name)
		if err != nil {
			if errors.Is(err, services.ErrInvalidName) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "nome é obrigatório"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível criar o ambiente"})
			return
		}
		c.JSON(http.StatusCreated, e)
	}
}

func ListEnvironments(svc *services.EnvironmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := userIDFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			return
		}
		list, err := svc.List(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível listar ambientes"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func UpdateEnvironment(svc *services.EnvironmentService) gin.HandlerFunc {
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
		var body environmentNameBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		e, err := svc.Update(userID, uri.ID, body.Name)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "ambiente não encontrado"})
				return
			}
			if errors.Is(err, services.ErrInvalidName) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "nome é obrigatório"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível atualizar o ambiente"})
			return
		}
		c.JSON(http.StatusOK, e)
	}
}

func DeleteEnvironment(svc *services.EnvironmentService) gin.HandlerFunc {
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
		if err := svc.Delete(userID, uri.ID); err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "ambiente não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível remover o ambiente"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
