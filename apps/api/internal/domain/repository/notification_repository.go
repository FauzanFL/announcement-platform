package repository

import (
	"announcement-api/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	CreateOne(ctx context.Context, notif *entity.Notification) error
	CreateBatch(ctx context.Context, notifs []entity.Notification) error
	FindByUser(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error)
	ExistsForUserAndAnnouncement(ctx context.Context, userID, announcementID uuid.UUID) (bool, error)
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) (int64, error)
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	DeleteByAnnouncementID(ctx context.Context, announcementID uuid.UUID) error
}
