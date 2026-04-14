package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"simon/apps/api/internal/config"
	"simon/apps/api/internal/middlewares"
	"simon/apps/api/internal/services"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	authSvc := services.NewAuthService(db, cfg.JWTSecret)
	colSvc := services.NewCollectionService(db)
	reqSvc := services.NewRequestService(db)

	r.GET("/health", Health(db))

	auth := r.Group("/auth")
	{
		auth.POST("/register", RegisterAuth(authSvc))
		auth.POST("/login", LoginAuth(authSvc))
	}

	protected := r.Group("", middlewares.AuthRequired(cfg.JWTSecret))
	{
		protected.GET("/me", Me(db))

		protected.POST("/collections", CreateCollection(colSvc))
		protected.GET("/collections", ListCollections(colSvc))
		protected.PUT("/collections/:id", UpdateCollection(colSvc))
		protected.DELETE("/collections/:id", DeleteCollection(colSvc))

		protected.POST("/collections/:id/requests", CreateRequest(reqSvc))
		protected.GET("/collections/:id/requests", ListRequests(reqSvc))
		protected.PUT("/requests/:id", UpdateRequest(reqSvc))
		protected.DELETE("/requests/:id", DeleteRequest(reqSvc))
	}
}
