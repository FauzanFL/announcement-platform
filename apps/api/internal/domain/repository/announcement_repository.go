package repository

import (
	"announcement-api/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type AnnouncementRepository interface {
	Create(ctx context.Context, announcement *entity.Announcement) error
	Update(ctx context.Context, announcement *entity.Announcement) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Announcement, error)
	FindAll(ctx context.Context) ([]entity.Announcement, error)
	FindAllWithReadStatus(ctx context.Context, userID uuid.UUID) ([]AnnouncementWithStatus, error)
	CountUnreadByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}

type AnnouncementWithStatus struct {
    entity.Announcement
    IsRead bool `json:"is_read"`
}