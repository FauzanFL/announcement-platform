package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	AnnouncementID uuid.UUID    `gorm:"type:uuid;not null;index" json:"announcement_id"`
	Announcement   Announcement `gorm:"foreignKey:AnnouncementID" json:"announcement"`
	IsRead         bool         `gorm:"default:false" json:"is_read"`
	CreatedAt      time.Time    `json:"created_at"`
	ReadAt         *time.Time   `json:"read_at"`
}
