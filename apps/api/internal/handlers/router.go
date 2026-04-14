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
	headerSvc := services.NewRequestHeaderService(db)
	querySvc := services.NewRequestQueryParamService(db)
	envSvc := services.NewEnvironmentService(db)
	envVarSvc := services.NewEnvironmentVariableService(db)

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

		protected.POST("/requests/:id/headers", CreateRequestHeader(headerSvc))
		protected.GET("/requests/:id/headers", ListRequestHeaders(headerSvc))
		protected.PUT("/headers/:id", UpdateRequestHeader(headerSvc))
		protected.DELETE("/headers/:id", DeleteRequestHeader(headerSvc))

		protected.POST("/requests/:id/query-params", CreateRequestQueryParam(querySvc))
		protected.GET("/requests/:id/query-params", ListRequestQueryParams(querySvc))
		protected.PUT("/query-params/:id", UpdateRequestQueryParam(querySvc))
		protected.DELETE("/query-params/:id", DeleteRequestQueryParam(querySvc))

		protected.POST("/environments", CreateEnvironment(envSvc))
		protected.GET("/environments", ListEnvironments(envSvc))
		protected.PUT("/environments/:id", UpdateEnvironment(envSvc))
		protected.DELETE("/environments/:id", DeleteEnvironment(envSvc))

		protected.POST("/environments/:id/variables", CreateEnvironmentVariable(envVarSvc))
		protected.GET("/environments/:id/variables", ListEnvironmentVariables(envVarSvc))
		protected.PUT("/environment-variables/:id", UpdateEnvironmentVariable(envVarSvc))
		protected.DELETE("/environment-variables/:id", DeleteEnvironmentVariable(envVarSvc))
	}
}
