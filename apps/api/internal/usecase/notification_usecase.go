package usecase

import (
	"announcement-api/internal/domain/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationUsecase struct {
	notifRepo repository.NotificationRepository
	annRepo repository.AnnouncementRepository
}

func NewNotificationUsecase(notifRepo repository.NotificationRepository, annRepo repository.AnnouncementRepository) *NotificationUsecase {
	return &NotificationUsecase{notifRepo: notifRepo, annRepo: annRepo}
}

func (u *NotificationUsecase) ListWithStatus(ctx context.Context, userID uuid.UUID) ([]repository.AnnouncementWithStatus, error) {
	return u.annRepo.FindAllWithReadStatus(ctx, userID)
}

func (u *NotificationUsecase) UnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return u.notifRepo.CountUnread(ctx, userID)
}

func (u *NotificationUsecase) MarkRead(ctx context.Context, announcementID uuid.UUID, userID uuid.UUID) error {
	return u.notifRepo.MarkRead(ctx, announcementID, userID)
}

func (u *NotificationUsecase) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return u.notifRepo.MarkAllRead(ctx, userID)
}
