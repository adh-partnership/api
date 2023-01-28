package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUIDBase struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;type:binary(16);not null;unique_index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UUIDBase) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
