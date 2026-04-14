package models

import "time"

type Request struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Method       string    `gorm:"size:16;not null" json:"method"`
	URL          string    `gorm:"type:text;not null" json:"url"`
	Body         string    `gorm:"type:text" json:"body"`
	CollectionID uint      `gorm:"not null;index" json:"collection_id"`
	CreatedAt    time.Time `json:"created_at"`

	Collection Collection `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
