package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/services"
)

type requestWriteBody struct {
	Name   string `json:"name" binding:"required,min=1,max=255"`
	Method string `json:"method" binding:"required"`
	URL    string `json:"url" binding:"required"`
	Body   string `json:"body"`
}

type requestIDURI struct {
	ID uint `uri:"id" binding:"required"`
}

func CreateRequest(svc *services.RequestService) gin.HandlerFunc {
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
		var body requestWriteBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		req, err := svc.Create(userID, uri.ID, body.Name, body.Method, body.URL, body.Body)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "collection não encontrada"})
				return
			}
			if errors.Is(err, services.ErrInvalidName) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "nome é obrigatório"})
				return
			}
			if errors.Is(err, services.ErrInvalidURL) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "url inválida"})
				return
			}
			if errors.Is(err, services.ErrInvalidMethod) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "método inválido (use GET, POST, PUT ou DELETE)"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível criar o request"})
			return
		}
		c.JSON(http.StatusCreated, req)
	}
}

func ListRequests(svc *services.RequestService) gin.HandlerFunc {
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
		list, err := svc.List(userID, uri.ID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "collection não encontrada"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível listar requests"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func UpdateRequest(svc *services.RequestService) gin.HandlerFunc {
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
		var body requestWriteBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		req, err := svc.Update(userID, uri.ID, body.Name, body.Method, body.URL, body.Body)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "request não encontrado"})
				return
			}
			if errors.Is(err, services.ErrInvalidName) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "nome é obrigatório"})
				return
			}
			if errors.Is(err, services.ErrInvalidURL) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "url inválida"})
				return
			}
			if errors.Is(err, services.ErrInvalidMethod) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "método inválido (use GET, POST, PUT ou DELETE)"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível atualizar o request"})
			return
		}
		c.JSON(http.StatusOK, req)
	}
}

func DeleteRequest(svc *services.RequestService) gin.HandlerFunc {
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
				c.JSON(http.StatusNotFound, gin.H{"error": "request não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível remover o request"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
