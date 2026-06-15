package entity

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	Creator   User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Notifications []Notification `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE;"`
}

type AnnouncementEvent struct {
	Type         string        `json:"type"`
	Announcement *Announcement `json:"anouncement,omitempty"`
	ID           string        `json:"id,omitempty"`
}
