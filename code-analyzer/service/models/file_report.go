package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JSONB map[string]interface{}

type FileReport struct {
	gorm.Model
	Id         uint32    `gorm:"type:int8;primary_key;"`
	ReportId   uuid.UUID `gorm:"type:uuid;not null;"`
	Path       string    `gorm:"type:varchar(100);not null;"`
	Findings   JSONB     `gorm:"type:jsonb" json:"findings"`
	QueuedAt   time.Time
	ScanningAt time.Time
	FinishAt   time.Time
}
