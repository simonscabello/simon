package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"simon/apps/api/internal/models"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("abrir conexão gorm: %w", err)
	}
	return db, nil
}

func Init(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("obter sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping no banco: %w", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Collection{}, &models.Request{}); err != nil {
		return fmt.Errorf("migrar schema: %w", err)
	}
	return nil
}
