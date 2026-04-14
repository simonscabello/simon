package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/services"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=200"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterAuth(svc *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		user, err := svc.Register(req.Name, req.Email, req.Password)
		if err != nil {
			if errors.Is(err, services.ErrEmailTaken) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email já cadastrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível concluir o cadastro"})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func LoginAuth(svc *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
			return
		}
		token, err := svc.Login(req.Email, req.Password)
		if err != nil {
			if errors.Is(err, services.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "email ou senha incorretos"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "não foi possível autenticar"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": token, "token_type": "Bearer", "expires_in": 900})
	}
}
