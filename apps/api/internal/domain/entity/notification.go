package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" example:"770e8400-e29b-41d4-a716-446655440002"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_user_announcement" json:"user_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	AnnouncementID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_user_announcement" json:"announcement_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ReadAt         time.Time   `json:"read_at" example:"2024-08-17T08:00:00Z"`

	User           User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	Announcement   Announcement `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE;" json:"-"`
}
