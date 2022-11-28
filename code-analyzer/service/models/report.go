package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	Id        uuid.UUID `gorm:"type:uuid;primary_key;"`
	RepoId    uint32    `gorm:"type:int8;not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
