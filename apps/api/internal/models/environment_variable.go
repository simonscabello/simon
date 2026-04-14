package models

import "time"

type EnvironmentVariable struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EnvironmentID uint      `gorm:"not null;uniqueIndex:idx_env_var_env_key" json:"environment_id"`
	Key           string    `gorm:"size:255;not null;uniqueIndex:idx_env_var_env_key" json:"key"`
	Value         string    `gorm:"type:text" json:"value"`
	Enabled       bool      `gorm:"not null;default:true" json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`

	Environment Environment `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
