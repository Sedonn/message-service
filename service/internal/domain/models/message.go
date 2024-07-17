package models

import (
	"time"
)

type Message struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	Content     string     `gorm:"column:content;size:256"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	ProcessedAt *time.Time `gorm:"column:processed_at;default:null"`
}
