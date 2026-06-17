package repository

import (
	"announcement-api/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	MarkRead(ctx context.Context, announcementID uuid.UUID, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	IsRead(ctx context.Context, announcementID uuid.UUID, userID uuid.UUID) (bool, error)
	ListReadIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]entity.Notification, error)
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
}
