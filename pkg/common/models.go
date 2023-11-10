package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        *string    `gorm:"primary_key;unique;type:uuid;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"autoUpdateTime:false;column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoCreateTime:true;column:updated_at;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedBy uuid.UUID  `gorm:"type:uuid;column:created_by;not null" json:"created_by"`
	UpdatedBy uuid.UUID  `gorm:"type:uuid;column:updated_by;not null" json:"updated_by"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == nil {
		id := uuid.New().String()
		base.ID = &id
	}
	return nil
}
