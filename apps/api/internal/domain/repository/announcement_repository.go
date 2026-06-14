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
}
