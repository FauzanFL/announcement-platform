package entity

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title     string    `gorm:"not null" json:"title" example:"National Holiday"`
	Content   string    `gorm:"type:text;not null" json:"content" example:"National Holiday on 17 August"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by" example:"660e8400-e29b-41d4-a716-446655440001"`
	Creator   User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt time.Time `json:"created_at" example:"2024-08-17T08:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-08-17T08:00:00Z"`

	Notifications []Notification `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE;"`
}

type AnnouncementEvent struct {
	Type         string        `json:"type" example:"created"`
	Announcement *Announcement `json:"anouncement,omitempty"`
	ID           string        `json:"id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
}
