package models

import "time"

type RequestQueryParam struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RequestID uint      `gorm:"not null;index" json:"request_id"`
	Key       string    `gorm:"size:255;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Enabled   bool      `gorm:"not null;default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`

	Request Request `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
