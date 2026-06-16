package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" example:"770e8400-e29b-41d4-a716-446655440002"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	AnnouncementID uuid.UUID    `gorm:"type:uuid;not null;index" json:"announcement_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Announcement   Announcement `gorm:"foreignKey:AnnouncementID" json:"announcement"`
	IsRead         bool         `gorm:"default:false" json:"is_read" example:"false"`
	CreatedAt      time.Time    `json:"created_at" example:"2024-08-17T08:00:00Z"`
	ReadAt         *time.Time   `json:"read_at" example:"2024-08-17T08:00:00Z"`
}
