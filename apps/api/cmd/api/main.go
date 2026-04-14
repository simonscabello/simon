package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"simon/apps/api/internal/config"
	"simon/apps/api/internal/database"
	"simon/apps/api/internal/handlers"
	"simon/apps/api/internal/middlewares"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET é obrigatório para autenticação")
	}

	db, err := database.Connect(cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("banco: %v", err)
	}
	if err := database.Init(db); err != nil {
		log.Fatalf("init banco: %v", err)
	}

	r := gin.Default()
	middlewares.Setup(r)
	handlers.RegisterRoutes(r, db, cfg)

	addr := ":" + cfg.HTTPPort
	log.Printf("API escutando em %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("servidor: %v", err)
	}
}
